package handshake

type HandshakeData struct {
	ProtocolString []byte
	InfoHash       []byte
	PeerId         []byte
}

func NewHandshake(protocolString, infoHash, peerId []byte) *HandshakeData {
	return &HandshakeData{
		ProtocolString: protocolString,
		InfoHash:       infoHash,
		PeerId:         peerId,
	}
}

// <protocol string length><protocol string><reserved><info_hash><peer_id>
func (h *HandshakeData) Serialize() []byte {
	protocolStringLen := len(h.ProtocolString)
	buff := make([]byte, 49+protocolStringLen)

	buff[0] = byte(protocolStringLen)

	currIndex := 1
	currIndex += copy(buff[currIndex:], h.ProtocolString)
	currIndex += copy(buff[currIndex:], make([]byte, 8))
	currIndex += copy(buff[currIndex:], h.InfoHash)
	copy(buff[currIndex:], h.PeerId)

	return buff
}
