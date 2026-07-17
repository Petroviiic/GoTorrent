package peer

import "net"

type PeerClient struct {
	Id         int
	Choked     bool
	Interested bool
	Bitfield   []byte
	Conn       net.Conn
	Manager    *Manager
}

func NewPeerClient(conn net.Conn) *PeerClient {
	client := &PeerClient{
		Choked:     true,
		Interested: false,
		Bitfield:   nil,
		Conn:       conn,
		Manager:    nil,
	}
	return client
}
