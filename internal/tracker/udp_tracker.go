package tracker

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"
	"time"
)

const MAGIC_CONSTANT = 0x41727101980
const CONNECT_ACTION_CONSTANT = 0
const ANNOUNCE_ACTION_CONSTANT = 1

func sendUDPRequest(trackerURL string, infoHash, peerID []byte, left string) ([]*Peer, error) {
	url, err := url.Parse(trackerURL)

	if err != nil {
		return nil, err
	}
	addr, err := net.ResolveUDPAddr("udp", url.Host)
	if err != nil {
		log.Println("Couldn’t resolve address:", err)
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer conn.Close()

	connected := false
	var connectionID uint64
	for retry := 0; retry <= 2; retry++ {
		// waitTime := time.Duration(15*(1<<retry)) * time.Second
		waitTime := time.Duration(2500*(1<<retry)) * time.Millisecond

		transactionID := rand.Uint32()

		message := prepareConnectRequest(int32(transactionID))

		if _, err = conn.Write(message); err != nil {
			log.Printf("Send failed: %v", err)
			return nil, err
		}

		conn.SetReadDeadline(time.Now().Add(waitTime))

		buffer := make([]byte, 16)

		if _, _, err := conn.ReadFromUDP(buffer); err != nil {
			log.Printf("Attempt %d failed (timeout/error): %v", retry+1, err)
			continue
		}

		ok, connId := parseConnectResponse(buffer, transactionID)

		if ok {
			connected = true
			connectionID = connId
			break
		}
	}

	if !connected {
		return nil, fmt.Errorf("failed to connect to tracker after retries")
	}
	for retry := 0; retry <= 2; retry++ {
		// waitTime := time.Duration(15*(1<<retry)) * time.Second
		waitTime := time.Duration(2500*(1<<retry)) * time.Millisecond

		transactionID := rand.Uint32()
		message := prepareAnnounceRequest(connectionID, transactionID, [20]byte(infoHash), [20]byte(peerID), 6881)

		if _, err = conn.Write(message); err != nil {
			log.Printf("Send failed: %v", err)
			return nil, err
		}
		conn.SetReadDeadline(time.Now().Add(waitTime))

		buffer := make([]byte, 4096)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Attempt %d failed (timeout/error): %v", retry+1, err)
			continue
		}

		peers, err := parseAnnonceResponse(buffer[:n], transactionID)
		if err != nil {
			log.Printf("Announce parse failed: %v", err)
			continue
		}

		return peers, nil
	}
	return nil, nil
}

func prepareConnectRequest(transactionId int32) []byte {
	req := make([]byte, 16)

	binary.BigEndian.PutUint64(req[:8], MAGIC_CONSTANT)
	binary.BigEndian.PutUint32(req[8:12], CONNECT_ACTION_CONSTANT)
	binary.BigEndian.PutUint32(req[12:16], uint32(transactionId))

	return req
}
func parseConnectResponse(buffer []byte, transactionID uint32) (bool, uint64) {
	if len(buffer) != 16 {
		log.Printf("Wrong response size: %d bytes\n", len(buffer))
		return false, 0
	}

	// 0       32-bit integer  action          0 // connect
	// 4       32-bit integer  transaction_id
	// 8       64-bit integer  connection_id
	// 16

	action := binary.BigEndian.Uint32(buffer[0:4])
	if action != CONNECT_ACTION_CONSTANT {
		log.Printf("Wrong action code: %d\n", action)
		return false, 0
	}

	gotTransactionID := binary.BigEndian.Uint32(buffer[4:8])
	if transactionID != gotTransactionID {
		log.Printf("Wrong transaction id: %d\n", gotTransactionID)
		return false, 0
	}

	return true, binary.BigEndian.Uint64(buffer[8:16])
}
func prepareAnnounceRequest(connectionId uint64, transactionId uint32, infoHash, peerID [20]byte, port uint16) []byte {
	req := make([]byte, 98)

	binary.BigEndian.PutUint64(req[:8], connectionId)
	binary.BigEndian.PutUint32(req[8:12], ANNOUNCE_ACTION_CONSTANT)
	binary.BigEndian.PutUint32(req[12:16], transactionId)
	copy(req[16:36], infoHash[:])
	copy(req[36:56], peerID[:])
	binary.BigEndian.PutUint64(req[56:64], 0)
	binary.BigEndian.PutUint64(req[64:72], 0)
	binary.BigEndian.PutUint64(req[72:80], 0)
	binary.BigEndian.PutUint32(req[80:84], 0)
	binary.BigEndian.PutUint32(req[84:88], 0)
	binary.BigEndian.PutUint32(req[88:92], 0)

	var numWant int32 = -1
	binary.BigEndian.PutUint32(req[92:96], uint32(numWant))
	binary.BigEndian.PutUint16(req[96:98], port)

	return req
}

func parseAnnonceResponse(buffer []byte, transactionID uint32) ([]*Peer, error) {
	if len(buffer) < 20 {
		return nil, fmt.Errorf("Short response size: %d bytes\n", len(buffer))
	}

	// 0           32-bit integer  action          1 // announce
	// 4           32-bit integer  transaction_id
	// 8           32-bit integer  interval
	// 12          32-bit integer  leechers
	// 16          32-bit integer  seeders
	// 20 + 6 * n  32-bit integer  IP address
	// 24 + 6 * n  16-bit integer  TCP port
	// 20 + 6 * N

	action := binary.BigEndian.Uint32(buffer[0:4])
	if action != ANNOUNCE_ACTION_CONSTANT {
		return nil, fmt.Errorf("Wrong action code: %d\n", action)
	}

	gotTransactionID := binary.BigEndian.Uint32(buffer[4:8])
	if transactionID != gotTransactionID {
		return nil, fmt.Errorf("Wrong transaction id: %d\n", gotTransactionID)
	}

	//TODO
	// interval := binary.BigEndian.Uint32(buffer[8:12])

	//for stats, gui
	// leechers := binary.BigEndian.Uint32(buffer[12:16])
	// seeders := binary.BigEndian.Uint32(buffer[16:20])

	// fmt.Println(interval, leechers, seeders)

	peers := DecodePeerList(buffer[20:])
	return peers, nil
}
