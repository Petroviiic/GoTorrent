package tracker

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/Petroviiic/GoTorrent/internal/bencode"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

func GetPeers(torrentData *bencode.TorrentFile, infoHash, peerID []byte) ([]*Peer, error) {
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

	// return sendRequest(torrentData.Announce, infoHash, peerID, left)     // this works like the original first version of my code, checks only the official tracker

	trackerURLs := []string{}
	if strings.HasPrefix(torrentData.Announce, "http") {
		trackerURLs = append(trackerURLs, torrentData.Announce)
	}
	for _, tier := range torrentData.AnnounceList {
		for _, link := range tier {
			if strings.HasPrefix(link, "http") {
				trackerURLs = append(trackerURLs, link)
			}
		}
	}
	uniquePeers := make(map[*Peer]bool)
	for _, trackerURL := range trackerURLs {
		newPeers, err := sendRequest(trackerURL, infoHash, peerID, left)

		if err != nil {
			continue
		}

		for _, val := range newPeers {
			uniquePeers[val] = true
		}
	}

	peers := []*Peer{}
	for p := range uniquePeers {
		peers = append(peers, p)
	}
	return peers, nil
}

func sendRequest(trackerURL string, infoHash, peerID []byte, left string) ([]*Peer, error) {
	req, err := http.NewRequest("GET", trackerURL, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
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

	if _, ok := res["peers"].([]byte); !ok {
		fmt.Println("Failed: The variable is not a string.")
		return nil, fmt.Errorf("type assertion failed. peers is []interface {}, not []uint8 ")
	}
	peers := DecodePeerList(res["peers"].([]byte))
	return peers, nil
}

func DecodePeerList(peers []byte) []*Peer {
	res := []*Peer{}

	if len(peers) < 6 {
		fmt.Println("invalid peers size")
		return res
	}

	for i := 0; i+6 <= len(peers); {
		peer := &Peer{}

		peer.IP = make(net.IP, 4)
		copy(peer.IP, peers[i:i+4])
		peer.Port = binary.BigEndian.Uint16(peers[i+4 : i+6])

		res = append(res, peer)
		i += 6
	}

	return res
}
