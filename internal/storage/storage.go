package storage

import (
	"bufio"
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
func (s *Storage) testfunc() {
	file := s.files[0]
	reader := bufio.NewReader(file.Handler)

	buffer := make([]byte, s.pieceLength)
	counter := 0
	bufferStart := 0
	for {
		n, err := reader.Read(buffer[bufferStart:])
		bufferStart = 0
		if err != nil {
			if err.Error() == io.EOF.Error() {
				fmt.Println("tu sam", n, counter)
				break
			}
			log.Fatalf("error reading file: %s", err)
		}
		counter++
	}
}
func (s *Storage) ScanDiskForDownloaded(hashedPieces []byte) error {

	// s.testfunc()

	// return nil
	fmt.Println()
	fmt.Println()
	fmt.Println()
	buffer := make([]byte, s.pieceLength)

	test := map[int]struct{}{}
	bufferStart := 0
	i := 0

	counter := 0
	for fileIdx, file := range s.files {
		fmt.Println("curr file info", file.Size, file.Handler.Name())
		reader := bufio.NewReader(file.Handler)

		for {
			n, err := reader.Read(buffer[bufferStart:])
			if err != nil {
				if err == io.EOF {
					fmt.Println("tu sam", i, n, counter, fileIdx)
					//return nil
					break
				}
				log.Fatalf("error reading file: %s", err)
			}
			fmt.Println(n+bufferStart == len(buffer), n, bufferStart, n+bufferStart, len(buffer))
			if n+bufferStart == len(buffer) || (n+bufferStart != len(buffer) && fileIdx == len(s.files)-1) {
				expected := hashedPieces[i : i+20]
				sum := sha1.Sum(buffer[:n+bufferStart])
				if !bytes.Equal(sum[:], expected) {
					fmt.Println("netacno", n, i/20, file.Handler.Name())
					test[i/20] = struct{}{}
				}
				i += 20
				bufferStart = 0
				counter++
				clear(buffer)
			} else {
				bufferStart += n
			}

		}
	}
	fmt.Println(test, len(test), counter, s.totalSize)
	return nil
}

func (s *Storage) AddToBuffer(data []byte, index int64) {
	s.Buffer <- DownloadedPieceData{
		data:       data,
		pieceIndex: index,
	}
}
