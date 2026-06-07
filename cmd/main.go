package main

import (
	"fmt"
	"os"

	"github.com/Petroviiic/GoTorrent/internal/bencode"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run ./... <path_to_torrent_file>")
		os.Exit(1)
	}

	path := os.Args[1]

	// fmt.Println(path)

	torrentFile, infoDictEncode, err := bencode.LoadAndDecode(path)
	if err != nil {
		fmt.Printf("Fatal: error %v", err)
		os.Exit(1)
	}

	_ = torrentFile
	_ = infoDictEncode
}
