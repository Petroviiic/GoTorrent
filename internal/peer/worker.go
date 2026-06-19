package peer

import (
	"fmt"
	"io"
	"sync"

	"github.com/Petroviiic/GoTorrent/internal/message"
)

type PieceOfWork struct {
}
type PieceOfResult struct {
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
			continue
		}

		fmt.Println("success", msg)
		if msg.ID == message.Bitfield {
			p.Bitfield = msg.Payload
		}
	}
}
