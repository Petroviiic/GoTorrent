package storage

import (
	"os"
	"path/filepath"
)

type FileInfo struct {
	Size        int64
	GlobalStart int64
	GlobalEnd   int64
	Handler     *os.File
}

type Storage struct {
	files       []FileInfo
	pieceLength int64
	totalSize   int64

	Buffer chan DownloadedPieceData
}

func NewStorage(fileInfo []File, pieceLength int64) (*Storage, error) {
	globalOffset := 0

	storage := &Storage{
		totalSize: 0,
		Buffer:    make(chan DownloadedPieceData, 50),
	}

	for _, file := range fileInfo {
		if err := os.MkdirAll(filepath.Dir(file.Path), 0755); err != nil {
			storage.Close()
			return nil, err
		}

		handler, err := os.OpenFile(file.Path, os.O_RDWR|os.O_CREATE, 0666)

		if err != nil {
			storage.Close()
			return nil, nil
		}

		if err := Preallocate(handler, 0, int64(file.Size)); err != nil {
			storage.Close()
			return nil, err
		}

		storage.files = append(storage.files, FileInfo{
			Size:        int64(file.Size),
			GlobalStart: int64(globalOffset),
			GlobalEnd:   int64(globalOffset + file.Size),
			Handler:     handler,
		})
		storage.totalSize += int64(file.Size)
		globalOffset += file.Size
	}

	storage.pieceLength = pieceLength
	return storage, nil
}

func (s *Storage) Close() {
	for _, file := range s.files {
		if file.Handler != nil {
			file.Handler.Close()
		}
	}
}

func (s *Storage) AddToBuffer(data []byte, index int64) {
	s.Buffer <- DownloadedPieceData{
		data:       data,
		pieceIndex: index,
	}
}
