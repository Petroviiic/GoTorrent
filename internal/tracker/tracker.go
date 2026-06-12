package tracker

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Petroviiic/GoTorrent/internal/bencode"
)

type Peer struct {
	ipAdress string
	port     string
}

func GetPeers(torrentData *bencode.TorrentFile, infoHash, peerID []byte) ([]*Peer, error) {
	req, err := http.NewRequest("GET", torrentData.Announce, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()

	left := ""
	if len(torrentData.Info.Files) == 0 {
		left = fmt.Sprintf("%d", torrentData.Info.Length)
	} else {
		leftNum := 0
		for _, file := range torrentData.Info.Files {
			leftNum += file.Length
		}

		left = fmt.Sprintf("%d", leftNum)
	}

	if left == "" {
		return nil, fmt.Errorf("something went wrong. 'left' is empty")
	}
	fmt.Println(left)
	return nil, nil

	params.Add("info_hash", string(infoHash))
	params.Add("peer_id", string(peerID))
	params.Add("port", "6881")
	params.Add("uploaded", "0")
	params.Add("downloaded", "0")
	params.Add("left", left)
	params.Add("compact", "1")

	req.URL.RawQuery = params.Encode()

	resp, err := http.Get(req.URL.String())
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	peers, err := decodePeerBody(body)

	if err != nil {
		return nil, err
	}

	return peers, nil
}

func decodePeerBody(body []byte) ([]*Peer, error) {
	decoder := bencode.NewDecoder(body)
	res, err := decoder.Decode(decoder.Buffer, 0)

	if err != nil {
		return nil, err
	}

	if _, ok := res["peers"]; !ok {
		return nil, fmt.Errorf("something went wrong. peers not present in response body")
	}

	peers := decodePeerList(res["peers"].([]byte))

	return peers, nil
}

func decodePeerList(peers []byte) []*Peer {
	res := []*Peer{}

	return res
}
