package main

import (
	"fmt"
	"os"

	"github.com/Petroviiic/GoTorrent/internal/bencode"
	"github.com/Petroviiic/GoTorrent/internal/peer"
	"github.com/Petroviiic/GoTorrent/internal/tracker"
	"github.com/Petroviiic/GoTorrent/internal/utils"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run ./... <path_to_torrent_file>")
		os.Exit(1)
	}

	path := os.Args[1]

	torrentFile, infoHash, err := bencode.LoadAndDecode(path)

	if err != nil {
		fmt.Printf("Fatal: error %v", err)
		os.Exit(1)
	}
	fmt.Println("torrent file successfully loaded")

	peerID := utils.GeneratePeerID([]byte("-GO0001-"))

	peers, err := tracker.GetPeers(torrentFile, infoHash, peerID)

	if err != nil {
		fmt.Printf("Fatal: error %v", err)
		os.Exit(1)
	}
	fmt.Println("peers successfully retrieved")

	peer.ConnectToPeers(peers, infoHash, peerID)
}
