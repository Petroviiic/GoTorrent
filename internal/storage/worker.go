package storage

import (
	"log"
	"sync"
)

type DownloadedPieceData struct {
	data       []byte
	pieceIndex int64
}

func (s *Storage) StartWorker(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	for piece := range s.Buffer {
		if err := s.WriteAt(piece.data, piece.pieceIndex); err != nil {
			log.Printf("storage worker error on piece %d: %v\n", piece.pieceIndex, err)
		}
	}
}
