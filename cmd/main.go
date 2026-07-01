package main

import (
	"fmt"
	"os"
	"sync"

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
	fmt.Println("torrent file successfully loaded", len(torrentFile.Info.Pieces), len(torrentFile.Info.Files), torrentFile.Info.PieceLength)

	peerID := utils.GeneratePeerID([]byte("-GO0001-"))

	peers, err := tracker.GetPeers(torrentFile, infoHash, peerID)

	if err != nil {
		fmt.Printf("Fatal: error %v", err)
		os.Exit(1)
	}
	fmt.Printf("%d peers successfully retrieved\n", len(peers))

	workers := peer.ConnectToPeers(peers, infoHash, peerID)

	fmt.Printf("connected to %v clients\n", len(workers))
	//workChannel := make(chan peer.PieceOfWork, 100)
	workManager := peer.NewManager(torrentFile.Info.Pieces, torrentFile.Info.PieceLength)
	var wg sync.WaitGroup
	for _, worker := range workers {
		worker.Manager = workManager

		wg.Add(1)
		go worker.StartWorker(&wg)
	}

	wg.Wait()
}
