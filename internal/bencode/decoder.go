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

	decoder := newDecoder(buffer)

	if _, err := decoder.decode(buffer, 0); err != nil {
		return err
	}

	return nil
}

func (d *Decoder) decode(buffer []byte, index int) (any, error) {
	for i := index; i < len(buffer); {
		switch b := buffer[i]; {
		case b == 'i':
			res, newIndex, err := d.decoders[b](i)

			if err != nil {
				break
			}

			i = newIndex
			fmt.Println(res)
		case b == 'l':
			res, newIndex, err := d.decoders[b](i)

			if err != nil {
				break
			}

			i = newIndex
			fmt.Println(res)
		case b == 'd':
		case b >= '0' && b <= '9':
			res, newIndex, err := d.decoders[b](i)

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

func (d *Decoder) decodeInt(index int) (any, int, error) {
	end := index
	for i := index; i < len(d.buffer); i++ {
		b := (d.buffer)[i]
		if b == 'e' {
			end = i
			break
		}
	}
	num, err := strconv.Atoi(string(d.buffer[index:end]))
	return num, end, err
}

func (d *Decoder) decodeString(index int) (any, int, error) {
	end := index
	for i := index; i < len(d.buffer); i++ {
		b := (d.buffer)[i]
		if b == ':' {
			end = i
			break
		}
	}
	num, err := strconv.Atoi(string(d.buffer[index:end]))

	if err != nil {
		return "", -1, err
	}

	return string(d.buffer[end+1 : end+1+num]), end + num + 2, nil
}

func (d *Decoder) decodeList(index int) (any, int, error) {
	res := []any{}

	end := index
	for i := index; i < len(d.buffer); i++ {
		b := (d.buffer)[i]

		item, newIndex, err := d.decoders[b](i)

		if err != nil {
			return nil, -1, err
		}

		res = append(res, item)
		i = newIndex

		if b == 'e' {
			end = i
			break
		}
	}

	return res, end, nil
}

func (d *Decoder) decodeDictionary(index int) (any, int, error) {
	res := map[any]any{}

	end := index
	isKey := true
	var lastKey any
	for i := index; i < len(d.buffer); i++ {
		b := (d.buffer)[i]

		item, newIndex, err := d.decoders[b](i)

		if err != nil {
			return nil, -1, err
		}

		if isKey {
			res[item] = false
			lastKey = item
			isKey = false
		} else {
			res[lastKey] = item
			lastKey = nil
			isKey = true
		}
		i = newIndex

		if b == 'e' {
			end = i
			break
		}
	}

	return res, end, nil
}
