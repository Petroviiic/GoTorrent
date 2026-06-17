package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

type MessageID uint8

const (
	Choke          MessageID = 0
	Unchoke        MessageID = 1
	Interested     MessageID = 2
	Not_interested MessageID = 3
	Have           MessageID = 4
	Bitfield       MessageID = 5
	Request        MessageID = 6
	Piece          MessageID = 7
	Cancel         MessageID = 8
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

func Deserialize(r io.Reader) (*Message, error) {
	sizeBuffer := make([]byte, 4)

	_, err := io.ReadFull(r, sizeBuffer)
	if err != nil {
		return nil, err
	}

	msgSize := binary.BigEndian.Uint32(sizeBuffer)

	if msgSize == 0 {
		fmt.Println("keep alive message")
		return &Message{}, nil
	}

	msgBuffer := make([]byte, msgSize)
	_, err = io.ReadFull(r, msgBuffer)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		ID:      MessageID(msgBuffer[0]),
		Payload: make([]byte, msgSize-1),
	}
	copy(msg.Payload, msgBuffer[1:])
	return msg, nil
}
