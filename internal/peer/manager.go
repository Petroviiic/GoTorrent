package peer

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/Petroviiic/GoTorrent/internal/storage"
)

type Manager struct {
	workChannel      chan PieceOfWork
	TotalPieces      int
	activeWorkers    int32
	DoneChannel      chan struct{}
	once             sync.Once
	mutex            sync.Mutex
	Storage          *storage.Storage
	DownloadedPieces map[int]struct{}
	CompletedPieces  int64
}

// TODO : dodaj rarest first. ovaj trenutni approach je dobar ako zelim preview tipa film neki pa pieces moraju jedan za drugim da dolaze
func NewManager(pieces []byte, pieceSize int, totalLength int, storage *storage.Storage, startPiece, endPiece int) *Manager {
	totalPieces := len(pieces) / 20

	if startPiece < 0 {
		startPiece = 0
	}
	if endPiece <= 0 || endPiece > totalPieces {
		endPiece = totalPieces
	}
	if startPiece >= endPiece {
		log.Fatalf("invalid piece range: start %d >= end %d", startPiece, endPiece)
	}

	neededPiecesCount := endPiece - startPiece

	manager := &Manager{
		workChannel:      make(chan PieceOfWork, totalPieces),
		TotalPieces:      neededPiecesCount,
		activeWorkers:    0,
		DoneChannel:      make(chan struct{}),
		Storage:          storage,
		DownloadedPieces: make(map[int]struct{}),
		CompletedPieces:  0,
	}
	// for i, j := 28000, 1400; i < len(pieces); j++ {
	for i, j := startPiece*20, startPiece; i < len(pieces) && j < endPiece; j++ {
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
			Index:       j,
			Hash:        hashCopy,
			Length:      currentPieceLength,
			TotalBlocks: (currentPieceLength + BLOCK_SIZE - 1) / BLOCK_SIZE,
		}

		manager.workChannel <- workPiece

		i += 20
	}
	return manager
}

func (m *Manager) AddNewEntry(index int, data []byte) {
	m.mutex.Lock()

	if _, ok := m.DownloadedPieces[index]; ok {
		m.mutex.Unlock()
		log.Printf("piece with index %v already exists in the storage", index)
		return
	}
	m.DownloadedPieces[index] = struct{}{}

	currentLen := len(m.DownloadedPieces)

	if currentLen == m.TotalPieces {
		fmt.Println("download done")
		m.once.Do(func() {
			close(m.DoneChannel)
		})
	}

	m.mutex.Unlock()

	m.Storage.AddToBuffer(data, int64(index))

	fmt.Printf("new entry index %v, storage len %v\n", index, currentLen)
}

func (m *Manager) WorkerStarted() {
	atomic.AddInt32(&m.activeWorkers, 1)
}

func (m *Manager) WorkerDone() {
	remaining := atomic.AddInt32(&m.activeWorkers, -1)

	if remaining == 0 {
		log.Println("all peers disconnected...")

		m.once.Do(func() {
			close(m.DoneChannel)
		})
	}
}
