package storage

import (
	"path/filepath"

	"github.com/Petroviiic/GoTorrent/internal/bencode"
)

type File struct {
	Path string
	Size int
}

func NewFile(path string, size int) File {
	return File{
		Path: path,
		Size: size,
	}
}

func GetFilesList(torrentInfo *bencode.InfoDict) []File {
	files := []File{}

	if len(torrentInfo.Files) == 0 {
		files = append(files, NewFile(torrentInfo.Name, torrentInfo.Length))
		return files
	}

	for _, fileInfo := range torrentInfo.Files {
		elem := append([]string{torrentInfo.Name}, fileInfo.Path...)

		files = append(files, NewFile(filepath.Join(elem...), fileInfo.Length))
	}

	return files
}
