package message

import "encoding/binary"

type MessageID uint8

const (
	choke      MessageID = 0
	unchoke    MessageID = 1
	interested MessageID = 2
	// 0 - choke
	// 1 - unchoke
	// 2 - interested
	// 3 - not interested
	// 4 - have
	// 5 - bitfield
	// 6 - request
	// 7 - piece
	// 8 - cancel
)

type Message struct {
	ID      MessageID
	Payload []byte
}

// <length prefix><message ID><payload>
func (m *Message) Serialize() []byte {
	size := len(m.Payload) + 1
	buff := make([]byte, size)

	binary.BigEndian.PutUint32(buff[:4], uint32(size))
	buff[4] = byte(m.ID)
	copy(buff[5:], m.Payload)

	return buff
}
