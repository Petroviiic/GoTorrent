package tracker

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

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
	trackerURLs = append(trackerURLs,
		"http://tracker.opentrackr.org:1337/announce",
		"http://open.tracker.cl:1337/announce",
		"http://www.torrentsnipe.info:2701/announce",
		"http://www.genesis-sp.org:2710/announce",
		"http://tracker2.dler.org:80/announce",
		"http://tracker.zhuqiy.dgj055.icu:80/announce",
		"http://tracker.xiaoduola.xyz:6969/announce",
		"http://tracker.sbsub.com:2710/announce",
		"http://tracker.renfei.net:8080/announce",
		"http://tracker.qu.ax:6969/announce",
		"http://tracker.mywaifu.best:6969/announce",
		"http://tracker.lintk.me:2710/announce",
		"http://tracker.ipv6tracker.org:80/announce",
		"http://tracker.dler.org:6969/announce",
		"http://tracker.bt4g.com:2095/announce",
		"http://tracker.bt-hash.com:80/announce",
		"http://tracker.bittor.pw:1337/announce",
		"http://tracker.23794.top:6969/announce",
		"http://tr.kxmp.cf:80/announce",
		"http://seeders-paradise.org:80/announce",
		"http://lucke.fenesisu.moe:6969/announce",
		"http://buny.uk:6969/announce",
		"http://bittorrent-tracker.e-n-c-r-y-p-t.net:1337/announce",
		"http://1337.abcvg.info:80/announce",
		"http://tracker.zhuqiy.com:80/announce",
		"http://tracker.waaa.moe:6969/announce",
		"http://tracker.privateseedbox.xyz:2710/announce",
		"http://tracker.nexusstream.eu:6969/announce",
		"http://tracker.dler.com:6969/announce",
		"http://tracker.dhitechnical.com:6969/announce",
	)
	uniquePeers := make(map[string]*Peer)
	for _, trackerURL := range trackerURLs {
		newPeers, err := sendRequest(trackerURL, infoHash, peerID, left)

		if err != nil {
			continue
		}

		for _, p := range newPeers {
			key := fmt.Sprintf("%s:%d", p.IP, p.Port)
			uniquePeers[key] = p
		}
	}

	peers := []*Peer{}
	for _, p := range uniquePeers {
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
	params.Add("numwant", "50")

	// encodedInfoHash := urlEncodeBytes(infoHash)
	// encodedPeerID := urlEncodeBytes(peerID)
	// req.URL.RawQuery = fmt.Sprintf(
	// 	"info_hash=%s&peer_id=%s&port=6881&uploaded=0&downloaded=0&left=%s&compact=1&numwant=50&event=started", encodedInfoHash, encodedPeerID, left,
	// )
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req.URL.RawQuery = params.Encode()
	req = req.WithContext(ctx)

	fmt.Println(req.URL.String())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4*1024*1024))

	//fmt.Printf("Sirovi odgovor od trackera (string): %s\n", string(body))
	// fmt.Printf("Sirovi odgovor (hex/bytes len): %d\n", len(body))
	peers, err := decodePeerBody(body)
	if err != nil {
		return nil, err
	}

	return peers, nil
}

// func urlEncodeBytes(b []byte) string {
// 	var buf strings.Builder
// 	for _, byteVal := range b {
// 		buf.WriteByte('%')
// 		buf.WriteString(fmt.Sprintf("%02X", byteVal))
// 	}
// 	return buf.String()
// }

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
		fmt.Println("invalid peers size", len(peers))
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
