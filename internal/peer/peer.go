package peer

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/Petroviiic/GoTorrent/internal/handshake"
	"github.com/Petroviiic/GoTorrent/internal/message"
	"github.com/Petroviiic/GoTorrent/internal/tracker"
)

type PeerClient struct {
	Conn       net.Conn
	Choked     bool
	Interested bool
	Bitfield   []byte
}

func NewPeerClient(conn net.Conn) *PeerClient {
	client := &PeerClient{
		Conn:       conn,
		Choked:     true,
		Interested: false,
		Bitfield:   nil,
	}
	return client
}
func (p *PeerClient) handlePeerClient() {
	defer p.Conn.Close()

	for {
		msg, err := message.Deserialize(p.Conn)

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("success", msg)
	}
}
func ConnectToPeers(peers []*tracker.Peer, infoHash []byte, peerID []byte) {
	ourHandshake := handshake.NewHandshake([]byte("BitTorrent protocol"), infoHash, peerID)
	ours := ourHandshake.Serialize()

	for _, peerInfo := range peers {
		address := fmt.Sprintf("%v:%d", peerInfo.IP, peerInfo.Port)
		conn, err := net.DialTimeout("tcp", address, 3*time.Second)

		if err != nil {
			fmt.Println(err)
			continue
		}

		_, err = conn.Write(ours)

		if err != nil {
			log.Println("Write error:", err)
			continue
		}

		buffer := make([]byte, len(ours))

		_, err = io.ReadFull(conn, buffer)
		if err != nil && err != io.EOF {
			log.Println("Read error:", err)
			continue
		}

		theirs := handshake.Deserialize(buffer)
		if theirs == nil {
			fmt.Println("theirs is nil")
			continue
		}
		if !handshake.AcceptHandshake(ourHandshake, theirs) {
			log.Println("Handshake not accepted, closing connection:", peerInfo.IP, peerInfo.Port)
			continue
		}

		peerClient := NewPeerClient(conn)

		go peerClient.handlePeerClient()
	}

}
