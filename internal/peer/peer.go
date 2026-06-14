package peer

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Petroviiic/GoTorrent/handshake"
	"github.com/Petroviiic/GoTorrent/internal/tracker"
)

func ConnectToPeers(peers []*tracker.Peer, infoHash []byte, peerID []byte) {
	ourHandshake := handshake.NewHandshake([]byte("BitTorrent protocol"), infoHash, peerID)
	our := ourHandshake.Serialize()

	for _, peerInfo := range peers {
		address := fmt.Sprintf("%v:%d", peerInfo.IP, peerInfo.Port)
		conn, err := net.DialTimeout("tcp", address, 3*time.Second)

		if err != nil {
			fmt.Println(err)
			continue
		}

		_, err = conn.Write(our)

		if err != nil {
			log.Println("Write error:", err)
			continue
		}
		fmt.Println(conn, "success")
	}

}
