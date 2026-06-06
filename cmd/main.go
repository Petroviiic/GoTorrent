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

	fmt.Println(path)

	torrentFile, err := bencode.LoadAndDecode(path)
	if err != nil {
		fmt.Printf("Fatal: error %v", err)
		os.Exit(1)
	}
	buffer, err := os.ReadFile(path)
	encoded, err := bencode.Encode(torrentFile)
	if err != nil {
		fmt.Printf("Fatal: error %v", err)
		os.Exit(1)
	}
	_ = torrentFile
	_ = buffer
	_ = encoded

	// 	for i := range encoded {
	// 		if encoded[i] != buffer[i] {
	// 			fmt.Println("kurcina")
	// 		}
	// 	}
	// 	// fmt.Println(encoded == buffer)
	// 	//fmt.Println(buffer)
}
