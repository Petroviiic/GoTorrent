package peer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/Petroviiic/GoTorrent/internal/message"
)

const BLOCK_SIZE = 16384

type PieceOfWork struct {
	Index  int
	Hash   []byte
	Length int
}
type PieceOfResult struct {
	Index      int
	Downloaded []byte
}

func (p *PeerClient) StartWorker(wg *sync.WaitGroup) {
	defer func() {
		p.Conn.Close()
		wg.Done()
	}()

	blocksArrivedCount := 0
	blocksArrived := []*PieceOfResult{}
	var currentPiece *PieceOfWork
	for {
		msg, err := message.Deserialize(p.Conn)

		if err != nil {
			if err != io.EOF {
				fmt.Println("error deserializing message:", err)
			}
			if currentPiece != nil {
				p.Manager.workChannel <- *currentPiece
			}
			return
		}

		fmt.Println(currentPiece == nil, !p.Choked, bytes.Equal(p.Bitfield, []byte{0}))
		if currentPiece == nil && !p.Choked && !bytes.Equal(p.Bitfield, []byte{0}) {
			currentPiece = p.getNextAvailablePiece()
			if currentPiece != nil {
				blocksArrived = make([]*PieceOfResult, currentPiece.Length/BLOCK_SIZE)
				blocksArrivedCount = 0
				p.sendRequests(currentPiece)
			}
		}

		fmt.Println("success", msg)
		switch msg.ID {
		case message.Choke:
			p.Choked = true

			if currentPiece != nil {
				p.Manager.workChannel <- *currentPiece
				blocksArrivedCount = 0
				blocksArrived = nil
			}
		case message.Unchoke:
			p.Choked = false
		case message.Interested:
			p.Interested = true
		case message.Not_interested:
			p.Interested = false
		case message.Have:
			// index := binary.BigEndian.Uint32(msg.Payload[1:])
			index := binary.BigEndian.Uint32(msg.Payload[0:])
			p.UpdatePiece(int(index))
		case message.Bitfield:
			p.Bitfield = msg.Payload
		case message.Request:
		case message.Piece:
			//dosao piece koji sam requestovao

			if currentPiece == nil {
				continue
			}
			pieceOfResult := &PieceOfResult{}

			if blocksArrivedCount < currentPiece.Length {
				blocksArrived[pieceOfResult.Index] = pieceOfResult
				blocksArrivedCount++
			}

			if blocksArrivedCount == currentPiece.Length {
				if fullHash, ok := HashOk(blocksArrived, currentPiece.Hash); ok {
					//sacuvaj taj hash na disku, ili u mapi po indeksu currentpiece.Index
					p.Manager.AddNewEntry(currentPiece.Index, fullHash)
				} else {
					p.Manager.workChannel <- *currentPiece
				}

				// u svakom slucaju
				currentPiece = nil
				blocksArrived = nil

				if !p.Choked {
					currentPiece = p.getNextAvailablePiece()
					if currentPiece != nil {
						blocksArrived = make([]*PieceOfResult, currentPiece.Length/BLOCK_SIZE)
						blocksArrivedCount = 0
						p.sendRequests(currentPiece)
					}
				}
			}
		case message.Cancel:

		default:
			fmt.Println("unknown message type")
		}

		//pieceWork := <- workChannel

	}
}

func (p *PeerClient) getNextAvailablePiece() *PieceOfWork {
	fmt.Println("finding next available piece")
	for {
		select {
		case piece, ok := <-p.Manager.workChannel:
			if !ok {
				return nil
			}

			if p.HasPiece(piece.Index) {
				fmt.Println("next piece : ", piece)
				return &piece
			}
			p.Manager.workChannel <- piece
		default:
			return nil
		}

	}
}
func (p *PeerClient) sendRequests(currentPiece *PieceOfWork) {
	blocks := make([][]byte, currentPiece.Length/BLOCK_SIZE)
	fmt.Printf("sending requests for %v", currentPiece)
	for i := 0; i < len(blocks); i++ {
		if err := message.SendRequest(p.Conn, currentPiece.Index, i*BLOCK_SIZE, BLOCK_SIZE); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Printf("requests sent for %v", currentPiece)
}
