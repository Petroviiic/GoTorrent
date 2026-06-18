package network

import (
	"fmt"
	"net"

	"github.com/Petroviiic/GoTorrent/internal/message"
	"github.com/Petroviiic/GoTorrent/internal/peer"
)

type PieceOfWork struct {
}
type PieceOfResult struct {
}

type Worker struct {
	Conn      net.Conn
	PeerState *peer.PeerState
}

func NewWorker(conn net.Conn) *Worker {
	client := &Worker{
		Conn:      conn,
		PeerState: peer.NewPeerState(),
	}
	return client
}
func (p *Worker) StartWorker() {
	defer p.Conn.Close()

	for {
		msg, err := message.Deserialize(p.Conn)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if msg.ID == message.Bitfield {
			p.PeerState.Bitfield = msg.Payload
		}
		fmt.Println("success", msg)
	}
}
