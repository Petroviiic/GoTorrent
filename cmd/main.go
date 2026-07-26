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

const DOWNLOAD_PATH = "E:\\BittorrentClientTest"

// const DOWNLOAD_PATH = "C:\\BittorrentClientTest"

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

	//neka prvo provjeri da li neki pieces vec postoje u fajlu da ne skida ponovo dzaba sve
	//na disku je raw data, pa se mora provjeriti sha hash dijelova fajla!!

	//takodje mogu dodati procentnost skinutosti, kao pieces_in_queue/total_num_of_pieces
	storage, err := storage.NewStorage(fileList, int64(torrentFile.Info.PieceLength))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
	if len(workers) == 0 {
		return
	}
	workManager := peer.NewManager(torrentFile.Info.Pieces, torrentFile.Info.PieceLength, totalSize)
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	for i, worker := range workers {
		worker.Manager = workManager
		worker.Id = i + 1

		wg.Add(1)
		go worker.StartWorker(&wg, ctx)
	}

	<-workManager.DoneChannel
	cancel()

	wg.Wait()

	fmt.Println("done", len(workManager.Storage), len(workManager.Storage) == workManager.TotalPieces)

	if len(workManager.Storage) != workManager.TotalPieces {
		for i := range workManager.TotalPieces {
			if _, ok := workManager.Storage[i]; !ok { //&& i > 1400 {
				fmt.Println(i)
			}
		}
	}

	storage.Close()
}
