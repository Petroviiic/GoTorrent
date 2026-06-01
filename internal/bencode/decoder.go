package bencode

import (
	"fmt"
	"os"
	"strconv"
)

type TorrentFile struct {
	Announce     string     `bencode:"announce"`
	AnnounceList [][]string `bencode:"announce-list"`
	Comment      string     `bencode:"comment"`
	CreatedBy    string     `bencode:"created by"`
	CreationDate int        `bencode:"creation date"`
	Info         InfoDict   `bencode:"info"`
}

type InfoDict struct {
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

func LoadAndDecode(path string) error {
	//cekiraj jel validan path

	buffer, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	decoders := loadDecoders()

	if _, err := decoders.decode(buffer, 0); err != nil {
		return err
	}

	return nil
}

func (decoders *Decoders) decode(buffer []byte, index int) (any, error) {
	for i := index; i < len(buffer); {
		//switch b := buffer[i]; {
		// case b == 'i':
		// 	num, newIndex, err := decodeInt(buffer, i+1)

		// 	if err != nil {
		// 		break
		// 	}

		// 	i = newIndex
		// 	fmt.Println(num)
		// case b == 'l':
		// 	list, newIndex, err := decodeList(buffer, i)
		// 	if err != nil {
		// 		break
		// 	}

		// 	i = newIndex
		// 	fmt.Println(list)
		// case b == 'd':s
		// case b >= '0' && b <= '9':
		// 	text, newIndex, err := decodeString(buffer, i)
		// 	if err != nil {
		// 		break
		// 	}

		// 	i = newIndex
		// 	fmt.Println(text)
		// }

		switch b := buffer[i]; {
		case b == 'i':
			res, newIndex, err := (*decoders)[b](buffer, i)

			if err != nil {
				break
			}

			i = newIndex
			fmt.Println(res)
		case b == 'l':
			res, newIndex, err := (*decoders)[b](buffer, i)

			if err != nil {
				break
			}

			i = newIndex
			fmt.Println(res)
		case b == 'd':
		case b >= '0' && b <= '9':
			res, newIndex, err := (*decoders)[b](buffer, i)

			if err != nil {
				break
			}

			i = newIndex
			fmt.Println(res)
		}

		i++
	}
	return nil, nil
}

func decodeInt(buffer []byte, index int) (any, int, error) {
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

func decodeString(buffer []byte, index int) (any, int, error) {
	end := index
	for i := index; i < len(buffer); i++ {
		b := (buffer)[i]
		if b == ':' {
			end = i
			break
		}
	}

	num, err := strconv.Atoi(string(buffer[index:end]))

	if err != nil {
		return "", -1, err
	}

	return string(buffer[end+1 : end+1+num]), end + num + 2, nil
}

func decodeList(buffer []byte, index int) (any, int, error) {
	end := index
	for i := index; i < len(buffer); i++ {
		b := (buffer)[i]

		if b == 'e' {
			end = i
			break
		}
	}

	return []int{}, end, nil
}
