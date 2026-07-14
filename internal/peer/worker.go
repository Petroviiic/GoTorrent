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
const BLOCKS_SENT_PER_REQUEST = 5

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

	startBlockIndex := 0
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

		fmt.Println(currentPiece == nil, !p.Choked, p.Bitfield != nil, !bytes.Equal(p.Bitfield, []byte{0}))
		if currentPiece == nil && !p.Choked && p.Bitfield != nil && !bytes.Equal(p.Bitfield, []byte{0}) {
			currentPiece = p.getNextAvailablePiece()
			if currentPiece != nil {
				blocksArrived = make([]*PieceOfResult, currentPiece.Length/BLOCK_SIZE)
				blocksArrivedCount = 0
				startBlockIndex = 0
				p.sendRequests(currentPiece, startBlockIndex)
			}
		}

		fmt.Printf("success ")
		if msg.ID != message.Piece {
			fmt.Println(msg)
		} else {
			fmt.Println("new piece")
		}
		switch msg.ID {
		case message.Choke:
			p.Choked = true

			if currentPiece != nil {
				p.Manager.workChannel <- *currentPiece
				blocksArrivedCount = 0
				startBlockIndex = 0
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

			if currentPiece == nil || blocksArrived == nil {
				continue
			}
			pieceOfResult := &PieceOfResult{}

			if blocksArrivedCount < currentPiece.Length/BLOCK_SIZE {
				blocksArrived[pieceOfResult.Index] = pieceOfResult
				blocksArrivedCount++
			}
			startBlockIndex += BLOCKS_SENT_PER_REQUEST
			fmt.Println(blocksArrivedCount, currentPiece.Length/BLOCK_SIZE)
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
						startBlockIndex = 0
						p.sendRequests(currentPiece, startBlockIndex)
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
func (p *PeerClient) sendRequests(currentPiece *PieceOfWork, startBlockIndex int) {
	blocks := make([][]byte, currentPiece.Length/BLOCK_SIZE)
	fmt.Printf("sending requests for %v\n", currentPiece)

	endBlockIndex := startBlockIndex + BLOCKS_SENT_PER_REQUEST //ovo dodaj da salje 5 po 5 npr ako bude blokirao... a blokira. ugl dodaj mzd kao parametre funkcije pocetni i krajnji indeks koji se trebaju poslati
	if len(blocks) < endBlockIndex {
		endBlockIndex = len(blocks)
	}
	for i := startBlockIndex; i < endBlockIndex; i++ {
		if err := message.SendRequest(p.Conn, currentPiece.Index, i*BLOCK_SIZE, BLOCK_SIZE); err != nil {
			fmt.Println(err)
		}
	}

	fmt.Printf("requests sent for %v\n", currentPiece)
}
