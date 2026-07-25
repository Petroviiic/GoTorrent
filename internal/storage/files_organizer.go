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

func GetFilesList(downloadPath string, torrentInfo *bencode.InfoDict) ([]File, int) {
	files := []File{}

	if len(torrentInfo.Files) == 0 {
		files = append(files, NewFile(filepath.Join(downloadPath, torrentInfo.Name), torrentInfo.Length))
		return files, torrentInfo.Length
	}

	totalSize := 0
	for _, fileInfo := range torrentInfo.Files {
		elem := append([]string{downloadPath, torrentInfo.Name}, fileInfo.Path...)

		files = append(files, NewFile(filepath.Join(elem...), fileInfo.Length))
		totalSize += fileInfo.Length
	}

	return files, totalSize
}
