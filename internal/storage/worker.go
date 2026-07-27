package storage

import (
	"context"
	"log"
	"sync"
)

type DownloadedPieceData struct {
	data       []byte
	pieceIndex int64
}

func (s *Storage) StartWorker(wg *sync.WaitGroup, ctx context.Context) {
	defer func() {
		wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			return

		case piece, ok := <-s.Buffer:
			if !ok {
				continue
			}

			if err := s.WriteAt(piece.data, piece.pieceIndex); err != nil {
				log.Println("storage worker error : ", err)
				//s.buffer <- piece
				continue
			}
		}
	}
}
