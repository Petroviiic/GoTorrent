package peer

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/Petroviiic/GoTorrent/internal/message"
)

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

	for {
		msg, err := message.Deserialize(p.Conn)

		if err != nil {
			if err != io.EOF {
				fmt.Println("error deserializing message:", err)
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
			index := binary.BigEndian.Uint32(msg.Payload[1:])
			p.UpdatePiece(int(index))
		case message.Bitfield:
			p.Bitfield = msg.Payload
		case message.Request:

		case message.Piece:
			//dosao piece koji sam requestovao
		case message.Cancel:

		default:
			fmt.Println("unknown message type")
		}

		//pieceWork := <- workChannel

	}
}
