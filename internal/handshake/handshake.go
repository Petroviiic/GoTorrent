package handshake

import (
	"bytes"
	"fmt"
)

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

func AcceptHandshake(ours, theirs *HandshakeData) bool {
	return bytes.Equal(ours.InfoHash, theirs.InfoHash) && bytes.Equal(ours.ProtocolString, theirs.ProtocolString)
}
func Deserialize(buffer []byte) *HandshakeData {
	if len(buffer) < 68 {
		fmt.Println("handshake buffer too small")
		return nil
	}

	protocolStringLen := int(buffer[0])
	if protocolStringLen != 19 {
		fmt.Println("invalid protocol string length")
		return nil
	}

	curr := 1
	protocolString := make([]byte, protocolStringLen)
	curr += copy(protocolString, buffer[curr:curr+protocolStringLen])

	curr += 8 // 8 reserved bits
	infoHash := make([]byte, 20)
	curr += copy(infoHash, buffer[curr:curr+20])

	peerId := make([]byte, 20)
	copy(peerId, buffer[curr:curr+20])

	return NewHandshake(protocolString, infoHash, peerId)
}
