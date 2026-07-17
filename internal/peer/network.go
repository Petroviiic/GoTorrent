package peer

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Petroviiic/GoTorrent/internal/handshake"
	"github.com/Petroviiic/GoTorrent/internal/message"
	"github.com/Petroviiic/GoTorrent/internal/tracker"
)

func ConnectToPeers(peers []*tracker.Peer, infoHash []byte, peerID []byte) []*PeerClient {
	var mutex sync.Mutex

	connectedPeers := []*PeerClient{}

	ourHandshake := handshake.NewHandshake([]byte("BitTorrent protocol"), infoHash, peerID)
	ours := ourHandshake.Serialize()

	var wg sync.WaitGroup

	for i, peerInfo := range peers {
		wg.Add(1)

		go func() {
			defer wg.Done()

			address := fmt.Sprintf("%v:%d", peerInfo.IP, peerInfo.Port)
			conn, err := net.DialTimeout("tcp", address, 3*time.Second)

			if err != nil {
				fmt.Printf("%d: %v\n", i+1, err)
				return
			}
			connectionAlive := false

			defer func() {
				if !connectionAlive {
					conn.Close()
				}
			}()

			conn.SetDeadline(time.Now().Add(3 * time.Second))
			_, err = conn.Write(ours)

			if err != nil {
				log.Println("Write error:", err)
				return
			}

			buffer := make([]byte, len(ours))

			_, err = io.ReadFull(conn, buffer)
			if err != nil && err != io.EOF {
				log.Println("Read error:", err)
				return
			}

			theirs := handshake.Deserialize(buffer)
			if theirs == nil {
				fmt.Println("theirs is nil")
				return
			}
			if !handshake.AcceptHandshake(ourHandshake, theirs) {
				log.Println("Handshake not accepted, closing connection:", peerInfo.IP, peerInfo.Port)
				return
			}
			if err := message.SendInterested(conn); err != nil {
				fmt.Println("send interested message failed")
				return
			}

			conn.SetDeadline(time.Time{})

			peerClient := NewPeerClient(conn)

			mutex.Lock()
			connectedPeers = append(connectedPeers, peerClient)
			connectionAlive = true
			mutex.Unlock()
		}()

	}

	wg.Wait()
	return connectedPeers
}
