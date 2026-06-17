package message

import "encoding/binary"

type MessageID uint8

const (
	choke          MessageID = 0
	unchoke        MessageID = 1
	interested     MessageID = 2
	not_interested MessageID = 3
	have           MessageID = 4
	bitfield       MessageID = 5
	request        MessageID = 6
	piece          MessageID = 7
	cancel         MessageID = 8
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

func (m *Message) Deserialize(data []byte) *Message {
	msg := &Message{}
	if len(data) < 5 {
		return msg
	}

	msgLen := binary.BigEndian.Uint32(data[:4])
	msgId := uint8(data[4])

	msg.ID = MessageID(msgId)
	copy(msg.Payload, data[5:msgLen])

	return msg
}
