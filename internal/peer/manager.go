package peer

import (
	"fmt"
	"log"
	"sync"

	"github.com/Petroviiic/GoTorrent/internal/storage"
)

type Manager struct {
	workChannel      chan PieceOfWork
	TotalPieces      int
	DoneChannel      chan struct{}
	mutex            sync.Mutex
	Storage          *storage.Storage
	DownloadedPieces map[int]struct{}
	CompletedPieces  int64
}

// TODO : dodaj rarest first. ovaj trenutni approach je dobar ako zelim preview tipa film neki pa pieces moraju jedan za drugim da dolaze
func NewManager(pieces []byte, pieceSize int, totalLength int, storage *storage.Storage) *Manager {
	totalPieces := len(pieces) / 20

	manager := &Manager{
		workChannel:      make(chan PieceOfWork, totalPieces),
		TotalPieces:      totalPieces,
		DoneChannel:      make(chan struct{}),
		Storage:          storage,
		DownloadedPieces: make(map[int]struct{}),
		CompletedPieces:  0,
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

func (m *Manager) AddNewEntry(index int, data []byte) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.DownloadedPieces[index]; ok {
		log.Printf("piece with index %v already exists in the storage", index)
		return
	}
	m.DownloadedPieces[index] = struct{}{}

	m.Storage.AddToBuffer(data, int64(index))

	fmt.Printf("new entry index %v, storage len %v\n", index, len(m.DownloadedPieces))

	if len(m.DownloadedPieces) == m.TotalPieces {
		// if len(m.Storage) == 22 {
		fmt.Println("download done")
		close(m.DoneChannel)
	}
}
