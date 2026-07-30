package storage

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Size        int64 //file size in bytes
	GlobalStart int64 //begining of storage (and all files of course) is 0
	GlobalEnd   int64 //basically GlobalStart+Size
	Handler     *os.File
}

type Storage struct {
	files       []FileInfo
	pieceLength int64 //in bytes
	totalSize   int64 //all files size

	Buffer     chan DownloadedPieceData
	Downloaded map[int]struct{}
}

func NewStorage(fileInfo []File, pieceLength int64) (*Storage, error) {
	globalOffset := 0

	storage := &Storage{
		totalSize:  0,
		Buffer:     make(chan DownloadedPieceData, 50),
		Downloaded: make(map[int]struct{}),
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

func (s *Storage) ScanDiskForDownloaded(hashedPieces []byte) error {
	fmt.Println("scanning disk for downloaded pieces...")
	buffer := make([]byte, s.pieceLength)

	bufferStart := 0
	i := 0

	counter := 0
	for fileIdx, file := range s.files {
		for {
			n, err := io.ReadFull(file.Handler, buffer[bufferStart:])

			totalRead := bufferStart + n
			if totalRead == len(buffer) || (fileIdx == len(s.files)-1 && (err == io.EOF || err == io.ErrUnexpectedEOF)) {
				expected := hashedPieces[i : i+20]
				sum := sha1.Sum(buffer[:totalRead])

				if bytes.Equal(sum[:], expected) {
					s.Downloaded[i/20] = struct{}{}
				}

				i += 20
				bufferStart = 0
				counter++
			} else {
				bufferStart = totalRead
			}

			if err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					break
				}
				log.Printf("error reading file: %s\n", err)
				return err
			}
		}
	}

	return nil
}

func (s *Storage) AddToBuffer(data []byte, index int64) {
	s.Buffer <- DownloadedPieceData{
		data:       data,
		pieceIndex: index,
	}
}
