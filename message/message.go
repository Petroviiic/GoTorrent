package message

type MessageID int

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

func (m *Message) Serialize() []byte {
	buff := make([]byte, 0)

	return buff
}
