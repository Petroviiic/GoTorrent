package peer

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
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
	PieceIndex  int
	BlockOffset int
	Downloaded  []byte
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
				log.Printf("peer %v error deserializing message: %v\n", p.Id, err)
			}
			if currentPiece != nil {
				p.Manager.workChannel <- *currentPiece
			}
			return
		}

		if currentPiece == nil && len(p.Manager.workChannel) == 0 {
			log.Printf("peer %v finished", p.Id)
			return
		}
		log.Printf("peer %v %v %v %v %v\n", p.Id, currentPiece == nil, !p.Choked, p.Bitfield != nil, !bytes.Equal(p.Bitfield, []byte{0}))
		if currentPiece == nil && !p.Choked && p.Bitfield != nil && !bytes.Equal(p.Bitfield, []byte{0}) {
			currentPiece = p.getNextAvailablePiece()
			if currentPiece != nil {
				blocksArrived = make([]*PieceOfResult, currentPiece.Length/BLOCK_SIZE)
				blocksArrivedCount = 0
				startBlockIndex = 0
				p.sendRequests(currentPiece, startBlockIndex)
			}
		}

		if msg.ID != message.Piece {
			log.Printf("peer %v success %v\n", p.Id, msg)
		} else {
			log.Printf("peer %v success new piece\n", p.Id)
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
				log.Printf("peer %v piece received but skipping it because : %v %v\n", p.Id, currentPiece == nil, blocksArrived == nil)
				continue
			}
			pieceOfResult := DecodePiece(msg.Payload)
			//log.Println(pieceOfResult.PieceIndex, pieceOfResult.BlockOffset/BLOCK_SIZE)

			if blocksArrivedCount < currentPiece.Length/BLOCK_SIZE {
				blocksArrived[pieceOfResult.BlockOffset/BLOCK_SIZE] = pieceOfResult
				blocksArrivedCount++
				//log.Println("blocks arrived ", blocksArrived)

				if blocksArrivedCount%BLOCKS_SENT_PER_REQUEST == 0 {
					startBlockIndex += BLOCKS_SENT_PER_REQUEST
					p.sendRequests(currentPiece, startBlockIndex)
				}
			}
			log.Println(blocksArrivedCount, currentPiece.Length/BLOCK_SIZE, startBlockIndex)
			if blocksArrivedCount == currentPiece.Length/BLOCK_SIZE {
				if fullHash, ok := HashOk(blocksArrived, currentPiece.Hash); ok {
					//sacuvaj taj hash na disku, ili u mapi po indeksu currentpiece.Index
					p.Manager.AddNewEntry(currentPiece.Index, fullHash)
				} else {
					log.Printf("peer %v wrong hash for piece %v\n", p.Id, currentPiece.Index)
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
			log.Printf("peer %v unknown message type\n", p.Id)
		}

		//pieceWork := <- workChannel

	}
}

func (p *PeerClient) getNextAvailablePiece() *PieceOfWork {
	log.Printf("peer %v finding next available piece\n", p.Id)
	log.Printf("peer %d bitfield: %v, length of workchannel %v\n", p.Id, p.Bitfield, len(p.Manager.workChannel))
	i := 0
	for {
		if i >= p.Manager.TotalPieces {
			log.Printf("peer %v couldnt find any available pieces\n", p.Id)
			return nil
		}
		select {
		case piece, ok := <-p.Manager.workChannel:
			if !ok {
				return nil
			}

			if p.HasPiece(piece.Index) {
				log.Printf("peer %v next piece : %v\n", p.Id, piece)
				return &piece
			}

			i++
			p.Manager.workChannel <- piece
		default:
			return nil
		}

	}
}
func (p *PeerClient) sendRequests(currentPiece *PieceOfWork, startBlockIndex int) {
	blocks := make([][]byte, currentPiece.Length/BLOCK_SIZE)
	log.Printf("peer %d sending requests for %v\n", p.Id, currentPiece)

	endBlockIndex := startBlockIndex + BLOCKS_SENT_PER_REQUEST //ovo dodaj da salje 5 po 5 npr ako bude blokirao... a blokira. ugl dodaj mzd kao parametre funkcije pocetni i krajnji indeks koji se trebaju poslati
	if len(blocks) < endBlockIndex {
		endBlockIndex = len(blocks)
	}
	for i := startBlockIndex; i < endBlockIndex; i++ {
		if err := message.SendRequest(p.Conn, currentPiece.Index, i*BLOCK_SIZE, BLOCK_SIZE); err != nil {
			log.Println(err)
		}
	}

	log.Printf("peer %d requests sent for %v\n", p.Id, currentPiece)
}
