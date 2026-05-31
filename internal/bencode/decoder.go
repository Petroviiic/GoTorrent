package bencode

import (
	"fmt"
	"os"
	"strconv"
)

type TorrentFile struct {
	Announce string   `bencode:"announce"`
	Info     InfoDict `bencode:"info"`
}

type InfoDict struct {
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
}

func Decode(path string) error {
	//cekiraj jel validan path

	buffer, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	for i := 0; i < len(buffer); {

		switch b := buffer[i]; {
		case b == 'i':
			num, newIndex, err := decodeInt(buffer, i+1)

			if err != nil {
				break
			}

			i = newIndex
			fmt.Println(num)
		case b == 'l':
		case b == 'd':
		case b >= '0' && b <= '9':
			num, newIndex, err := decodeString(buffer, i)
			if err != nil {
				break
			}

			i = newIndex
			fmt.Println(num)
		}
		i++
	}
	return nil
}

func decodeInt(buffer []byte, index int) (int, int, error) {
	end := index
	for i := index; i < len(buffer); i++ {
		b := (buffer)[i]
		if b == 'e' {
			end = i
			break
		}
	}
	// fmt.Println(string(buffer[index-1 : end]))
	num, err := strconv.Atoi(string(buffer[index:end]))
	// fmt.Println(num, err)

	return num, end, err
}
