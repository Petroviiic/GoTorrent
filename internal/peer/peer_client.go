package peer

import "net"

type PeerClient struct {
	Choked     bool
	Interested bool
	Bitfield   []byte
	Conn       net.Conn
	Manager    *Manager

	CurrentPiece *PieceOfWork
}

func NewPeerClient(conn net.Conn) *PeerClient {
	client := &PeerClient{
		Choked:       true,
		Interested:   false,
		Bitfield:     nil,
		Conn:         conn,
		Manager:      nil,
		CurrentPiece: nil,
	}
	return client
}
