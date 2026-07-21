package peer

import (
	"fmt"
	"log"
	"sync"
)

type Manager struct {
	workChannel chan PieceOfWork
	TotalPieces int
	Storage     map[int][]byte
	DoneChannel chan struct{}
	mutex       sync.Mutex
}

// TODO : dodaj rarest first. ovaj trenutni approach je dobar ako zelim preview tipa film neki pa pieces moraju jedan za drugim da dolaze
func NewManager(pieces []byte, pieceSize int, totalLength int) *Manager {
	totalPieces := len(pieces) / 20

	manager := &Manager{
		workChannel: make(chan PieceOfWork, totalPieces),
		Storage:     make(map[int][]byte),
		TotalPieces: totalPieces,
		DoneChannel: make(chan struct{}),
	}
	// for i, j := 28000, 1400; i < len(pieces); j++ {
	for i, j := 0, 0; i < len(pieces); j++ {
		endIndex := i + 20
		if endIndex > len(pieces) {
			endIndex = len(pieces)
		}
		hashCopy := make([]byte, endIndex-i)
		copy(hashCopy, pieces[i:endIndex])

		currentPieceLength := pieceSize
		if j == totalPieces-1 {
			currentPieceLength = totalLength - (j * pieceSize)
		}

		workPiece := PieceOfWork{
			Index:  j,
			Hash:   hashCopy,
			Length: currentPieceLength,
		}

		manager.workChannel <- workPiece

		i += 20
	}
	return manager
}

func (m *Manager) AddNewEntry(index int, hash []byte) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.Storage[index]; ok {
		log.Printf("piece with index %v already exists in the storage", index)
		return
	}
	m.Storage[index] = hash

	fmt.Printf("new entry index %v, storage len %v\n", index, len(m.Storage))

	if len(m.Storage) == m.TotalPieces {
		// if len(m.Storage) == 22 {
		fmt.Println("download done")
		close(m.DoneChannel)
	}
}
