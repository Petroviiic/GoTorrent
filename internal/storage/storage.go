package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Size    int64
	Start   int64
	Offset  int64
	Handler *os.File
}

type Storage struct {
	files     []FileInfo
	totalSize int64
}

func NewStorage(fileInfo []File) *Storage {
	globalOffset := 0

	storage := &Storage{
		totalSize: 0,
	}

	for _, file := range fileInfo {
		if err := os.MkdirAll(filepath.Dir(file.Path), 0755); err != nil {
			return nil
		}

		handler, err := os.OpenFile(file.Path, os.O_RDWR|os.O_CREATE, 0666)

		if err != nil {
			storage.Close()
			return nil
		}

		if err := Preallocate(handler, 0, int64(file.Size)); err != nil {
			fmt.Println(err)
			storage.Close()
			return nil
		}

		storage.files = append(storage.files, FileInfo{
			Size:    int64(file.Size),
			Start:   int64(globalOffset),
			Offset:  int64(globalOffset + file.Size),
			Handler: handler,
		})
		storage.totalSize += int64(file.Size)
		globalOffset += file.Size
	}

	return storage
}

func (s *Storage) Close() {
	for _, file := range s.files {
		file.Handler.Close()
	}
}
