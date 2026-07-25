//go:build windows

package storage

import (
	"os"
)

func Preallocate(file *os.File, offset, size int64) error {
	return file.Truncate(size + offset)
}
