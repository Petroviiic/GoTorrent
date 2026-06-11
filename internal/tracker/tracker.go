package tracker

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Petroviiic/GoTorrent/internal/bencode"
)

func GetPeers(torrentData *bencode.TorrentFile, infoHash, peerID []byte) {
	req, err := http.NewRequest("GET", torrentData.Announce, nil)
	if err != nil {
		log.Print(err)
		return
	}

	params := req.URL.Query()

	left := strconv.Itoa(torrentData.Info.Length)
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

	decodePeerBody(body)
}

func decodePeerBody(body []byte) {
	decoder := bencode.NewDecoder(body)
	res, err := decoder.Decode(decoder.Buffer, 0)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println(string(res["peers"].([]byte)))
}

func decodePeerList() {

}
