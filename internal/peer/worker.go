package peer

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/Petroviiic/GoTorrent/internal/message"
)

const BLOCK_SIZE = 16384
const BLOCKS_SENT_PER_REQUEST = 5
const MAX_WRONG_HASHES_ALLOWED = 6

type PieceOfWork struct {
	Index       int
	TotalBlocks int
	Hash        []byte
	Length      int
}
type PieceOfResult struct {
	PieceIndex  int
	BlockOffset int
	Downloaded  []byte
}

func (p *PeerClient) StartWorker(wg *sync.WaitGroup, ctx context.Context) {
	defer func() {
		p.Conn.Close()
		wg.Done()
		p.Manager.WorkerDone()
	}()

	p.Manager.WorkerStarted()

	startBlockIndex := 0
	blocksArrivedCount := 0
	blocksArrived := []*PieceOfResult{}
	var currentPiece *PieceOfWork

	var retrieveAndRequestPiece func()
	retrieveAndRequestPiece = func() {
		currentPiece = p.getNextAvailablePiece(ctx)
		if currentPiece != nil {
			blocksArrived = make([]*PieceOfResult, currentPiece.TotalBlocks)
			blocksArrivedCount = 0
			startBlockIndex = 0
			p.sendRequests(currentPiece, startBlockIndex)
		}
	}
	for {
		select {
		case <-ctx.Done():
			if currentPiece != nil {
				p.Manager.workChannel <- *currentPiece
			}
			return
		case <-p.Manager.DoneChannel:
			if currentPiece != nil {
				p.Manager.workChannel <- *currentPiece
			}
			return
		default:

			p.Conn.SetReadDeadline(time.Now().Add(120 * time.Second))
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

			switch msg.ID {
			case message.Choke:
				p.Choked = true

				if currentPiece != nil {
					p.Manager.workChannel <- *currentPiece
					blocksArrivedCount = 0
					startBlockIndex = 0
					currentPiece = nil
					blocksArrived = nil
				}
			case message.Unchoke:
				p.Choked = false
				if currentPiece == nil {
					retrieveAndRequestPiece()
				}
			case message.Interested:
				p.Interested = true
			case message.Not_interested:
				p.Interested = false
			case message.Have:
				index := binary.BigEndian.Uint32(msg.Payload[0:])
				p.UpdatePiece(int(index))

				if currentPiece == nil && !p.Choked {
					retrieveAndRequestPiece()
				}
			case message.Bitfield:
				p.Bitfield = msg.Payload
			case message.Request:
			case message.Piece:
				//dosao piece koji sam requestovao

				if currentPiece == nil || blocksArrived == nil {
					log.Printf("peer %v piece received but skipping it because : %v %v\n", p.Id, currentPiece == nil, blocksArrived == nil)

					if currentPiece != nil {
						log.Println("skipped ", currentPiece.Index)
					}
					continue
				}
				pieceOfResult := DecodePiece(msg.Payload)

				if blocksArrivedCount < currentPiece.TotalBlocks {
					if blocksArrived[pieceOfResult.BlockOffset/BLOCK_SIZE] == nil {
						blocksArrived[pieceOfResult.BlockOffset/BLOCK_SIZE] = pieceOfResult
						blocksArrivedCount++
					}

					if blocksArrivedCount%BLOCKS_SENT_PER_REQUEST == 0 {
						startBlockIndex += BLOCKS_SENT_PER_REQUEST
						p.sendRequests(currentPiece, startBlockIndex)
					}
				}

				if blocksArrivedCount == currentPiece.TotalBlocks {
					if downloadedData, ok := HashOk(blocksArrived, currentPiece.Hash); ok {
						p.Manager.AddNewEntry(currentPiece.Index, downloadedData)
					} else {
						fmt.Println("pogresan ", p.Id, p.Conn.RemoteAddr().String(), currentPiece, blocksArrivedCount, blocksArrived)
						log.Printf("peer %v wrong hash for piece %v\n", p.Id, currentPiece.Index)

						p.Manager.workChannel <- *currentPiece

						p.WrongHashes[currentPiece.Index]++
						p.TotalWrongHashes++

						if p.TotalWrongHashes >= MAX_WRONG_HASHES_ALLOWED {
							log.Printf("peer %v DISCONNECTED (too many bad hashes: %d)\n", p.Id, p.TotalWrongHashes)
							return
						}
					}

					currentPiece = nil
					blocksArrived = nil

					if !p.Choked {
						retrieveAndRequestPiece()
					}
				}
			case message.Cancel:

			default:
				log.Printf("peer %v unknown message type\n", p.Id)
			}
		}
	}
}

func (p *PeerClient) getNextAvailablePiece(ctx context.Context) *PieceOfWork {
	log.Printf("peer %v finding next available piece\n", p.Id)

	i := 0
	visited := make(map[int]struct{})
	for {
		if i >= p.Manager.TotalPieces {
			log.Printf("peer %v couldnt find any available pieces\n", p.Id)
			return nil
		}
		select {
		case <-ctx.Done():
			return nil
		case <-p.Manager.DoneChannel:
			return nil
		case piece, ok := <-p.Manager.workChannel:
			if !ok {
				return nil
			}

			i++
			if p.HasPiece(piece.Index) {
				if tries := p.WrongHashes[piece.Index]; tries >= 2 {
					log.Printf("peer %v skipping piece %d due to %d bad hash attempts\n", p.Id, piece.Index, tries)

					p.Manager.workChannel <- piece
					continue
				}
				log.Printf("peer %v next piece : %v\n", p.Id, piece)
				return &piece
			}

			if _, seen := visited[piece.Index]; seen {
				p.Manager.workChannel <- piece
				return nil
			}

			visited[piece.Index] = struct{}{}
			p.Manager.workChannel <- piece
		default:
			return nil
		}
	}
}
func (p *PeerClient) sendRequests(currentPiece *PieceOfWork, startBlockIndex int) {
	endBlockIndex := startBlockIndex + BLOCKS_SENT_PER_REQUEST
	if currentPiece.TotalBlocks < endBlockIndex {
		endBlockIndex = currentPiece.TotalBlocks
	}

	for i := startBlockIndex; i < endBlockIndex; i++ {
		reqLength := BLOCK_SIZE
		blockOffset := i * BLOCK_SIZE
		if blockOffset+reqLength > currentPiece.Length {
			reqLength = currentPiece.Length - blockOffset
		}

		if err := message.SendRequest(p.Conn, currentPiece.Index, blockOffset, reqLength); err != nil {
			log.Printf("ERROR SENDING REQUESTS; peer id: %v; error: %v\n", p.Id, err)
		}
	}
}
