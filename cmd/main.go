package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/Petroviiic/GoTorrent/internal/bencode"
	"github.com/Petroviiic/GoTorrent/internal/peer"
	"github.com/Petroviiic/GoTorrent/internal/storage"
	"github.com/Petroviiic/GoTorrent/internal/tracker"
	"github.com/Petroviiic/GoTorrent/internal/utils"
)

const DOWNLOAD_PATH = "E:\\GoBittorrentClient"

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

	fileList, totalSize := storage.GetFilesList(DOWNLOAD_PATH, &torrentFile.Info)
	fmt.Printf("torrent file successfully loaded; pieces count : %v, files %v, number of files %v, one piece length : %v\n", len(torrentFile.Info.Pieces)/20, fileList, len(fileList), torrentFile.Info.PieceLength)

	storage, err := storage.NewStorage(fileList, int64(torrentFile.Info.PieceLength))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := storage.ScanDiskForDownloaded(torrentFile.Info.Pieces); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	peerID := utils.GeneratePeerID([]byte("-GO0001-"))

	workManager := peer.NewManager(torrentFile.Info.Pieces, torrentFile.Info.PieceLength, totalSize, storage, 0, -1)

	if workManager.IsWorkChannelEmpty() {
		fmt.Println("all pieces are already downloaded")
		return
	}

	peers, err := tracker.GetPeers(torrentFile, infoHash, peerID)

	if err != nil {
		fmt.Printf("Fatal: error %v", err)
		os.Exit(1)
	}
	fmt.Printf("%d peers successfully retrieved\n", len(peers))

	workers := peer.ConnectToPeers(peers, infoHash, peerID)

	fmt.Printf("connected to %v clients\n", len(workers))

	if len(workers) == 0 {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())

	var storageWg sync.WaitGroup
	storageWg.Add(1)
	go storage.StartWorker(&storageWg)

	var peerWg sync.WaitGroup
	for i, worker := range workers {
		worker.Manager = workManager
		worker.Id = i + 1

		peerWg.Add(1)
		go worker.StartWorker(&peerWg, ctx)
	}

	<-workManager.DoneChannel
	cancel()

	peerWg.Wait()

	close(storage.Buffer)
	storageWg.Wait()

	fmt.Printf("download done: %v; number of downloaded pieces %d;\n", int(workManager.CompletedPieces) == workManager.TotalPieces, len(workManager.DownloadedPieces))

	storage.Close()
}
