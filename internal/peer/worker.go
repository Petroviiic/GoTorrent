package peer

import (
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
	blocksArrived := []PieceOfWork{}
	for {
		msg, err := message.Deserialize(p.Conn)

		if err != nil {
			if err != io.EOF {
				fmt.Println("error deserializing message:", err)
			}
			if p.CurrentPiece != nil {
				p.Manager.workChannel <- *p.CurrentPiece
			}
			return
		}

		fmt.Println("success", msg)
		switch msg.ID {
		case message.Choke:
			p.Choked = true
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
			blocksArrived = make([]PieceOfWork, p.CurrentPiece.Length)
		case message.Piece:
			blocksArrivedCount++
			if blocksArrivedCount == p.CurrentPiece.Length {

			} else {
				// pieceOfWork, err := message.RecievePiece()
				// if err != nil {
				// 	//uradi nesto
				// }

			}
			//dosao piece koji sam requestovao
		case message.Cancel:

		default:
			fmt.Println("unknown message type")
		}

		//pieceWork := <- workChannel

	}
}

func (p *PeerClient) getNextAvailablePiece() {
	nextPiece := <-p.Manager.workChannel

	if !p.HasPiece(nextPiece.Index) {
		p.Manager.workChannel <- nextPiece

	}
	p.CurrentPiece = &nextPiece

	fmt.Println("next piece : ", nextPiece)
	blocks := make([][]byte, nextPiece.Length/BLOCK_SIZE)
	for i := 0; i < len(blocks); i++ {
		if err := message.SendRequest(p.Conn, nextPiece.Index, i*BLOCK_SIZE, BLOCK_SIZE); err != nil {
			fmt.Println(err)
		}
	}
}
