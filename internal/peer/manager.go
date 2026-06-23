package peer

type Manager struct {
	workChannel chan PieceOfWork
}

// TODO : dodaj rarest first. ovaj trenutni approach je dobar ako zelim preview tipa film neki pa pieces moraju jedan za drugim da dolaze
func NewManager(pieces []byte, pieceSize int) *Manager {
	manager := &Manager{
		workChannel: make(chan PieceOfWork, len(pieces)/20),
	}

	for i, j := 0, 0; i < len(pieces); j++ {
		endIndex := i + 20
		if endIndex > len(pieces) {
			endIndex = len(pieces)
		}
		hashCopy := make([]byte, endIndex-i)
		copy(hashCopy, pieces[i:endIndex])

		workPiece := PieceOfWork{
			Index:  j,
			Hash:   hashCopy,
			Length: pieceSize,
		}

		manager.workChannel <- workPiece

		i += 20
	}
	return manager
}
