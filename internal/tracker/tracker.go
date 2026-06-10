package tracker

import (
	"log"
	"net/url"

	"github.com/Petroviiic/GoTorrent/internal/bencode"
)

func GetPeers(torrentData *bencode.TorrentFile, infoHash string, peerID string) {
	baseURL, err := url.Parse("https://example.com")
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}

	// http.Get(torrentData.)
}
