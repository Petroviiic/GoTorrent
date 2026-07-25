//go:build linux

package storage

import (
	"os"

	"golang.org/x/sys/unix"
)

func Preallocate(file *os.File, offset, size int64) error {
	return unix.Fallocate(int(file.Fd()), 0, offset, size)
}
