package peer

type PeerState struct {
	Choked     bool
	Interested bool
	Bitfield   []byte
}

func NewPeerState() *PeerState {
	client := &PeerState{
		Choked:     true,
		Interested: false,
		Bitfield:   nil,
	}
	return client
}
