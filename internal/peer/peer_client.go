package peer

import "net"

type PeerClient struct {
	Choked     bool
	Interested bool
	Bitfield   []byte
	Conn       net.Conn
}

func NewPeerClient(conn net.Conn) *PeerClient {
	client := &PeerClient{
		Choked:     true,
		Interested: false,
		Bitfield:   nil,
		Conn:       conn,
	}
	return client
}
