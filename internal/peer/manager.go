package peer

import (
	"fmt"
	"log"
)

type Manager struct {
	workChannel chan PieceOfWork
	TotalPieces int
	Storage     map[int][]byte
}

// TODO : dodaj rarest first. ovaj trenutni approach je dobar ako zelim preview tipa film neki pa pieces moraju jedan za drugim da dolaze
func NewManager(pieces []byte, pieceSize int) *Manager {
	totalPieces := len(pieces) / 20
	manager := &Manager{
		workChannel: make(chan PieceOfWork, totalPieces),
		Storage:     make(map[int][]byte),
		TotalPieces: totalPieces,
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

func (m *Manager) AddNewEntry(index int, hash []byte) {
	if _, ok := m.Storage[index]; ok {
		log.Printf("piece with index %v already exists in the storage", index)
		return
	}
	m.Storage[index] = hash

	fmt.Println("storage ", len(m.Storage))
}
