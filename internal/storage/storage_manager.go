package storage

import (
	"errors"
	"fmt"
)

func (s *Storage) WriteAt(data []byte, pieceIndex int64) error {
	globalPosition := pieceIndex * s.pieceLength

	if globalPosition+int64(len(data)) > s.totalSize {
		return errors.New("piece index greater than total torrent size")
	}
	if globalPosition < 0 {
		return errors.New("piece index cannot be negative")
	}

	bytesLeft := int64(len(data))
	currentByte := int64(0)

	for _, f := range s.files {
		if bytesLeft == 0 {
			break
		}

		currentGlobalPosition := globalPosition + currentByte

		if f.GlobalStart <= currentGlobalPosition && f.GlobalEnd > currentGlobalPosition {
			if f.Handler == nil {
				return errors.New("something went wrong. no files require the provided piece")
			}

			fileOffset := currentGlobalPosition - f.GlobalStart

			amountOfBytes := bytesLeft
			if currentGlobalPosition+bytesLeft > f.GlobalEnd {
				amountOfBytes = f.GlobalEnd - currentGlobalPosition
			}
			_, err := f.Handler.WriteAt(data[currentByte:currentByte+amountOfBytes], fileOffset)
			if err != nil {
				return err
			}

			currentByte += amountOfBytes
			bytesLeft -= amountOfBytes
		}
	}

	if bytesLeft > 0 {
		return fmt.Errorf("error : couldnt write all bytes")
	}
	return nil
}
