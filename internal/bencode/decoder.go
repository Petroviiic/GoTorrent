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

	decoder := NewDecoder(buffer)

	data, err := decoder.Decode(buffer, 0)
	if err != nil {
		return err
	}
	fmt.Println(data)
	return nil
}

func (d *Decoder) Decode(buffer []byte, index int) (any, error) {

	mainMap := map[any]any{}
	for i := index; i < len(buffer); {
		switch b := buffer[i]; {
		case b == 'i':
			res, newIndex, err := d.Decoders[b](i)

			if err != nil {
				//ovdje vjv treba return, jer ce biti infinity loop ako se desi greska
				break
			}

			i = newIndex
			fmt.Println(res)

			if res == 497360 {
				fmt.Println("EVO OME SPASAVAM TE OD SPAMA PLIZ SE SAUSTAVI")
				return nil, nil
			}
		case b == 'l':
			res, newIndex, err := d.Decoders[b](i)

			if err != nil {
				break
			}

			i = newIndex
			fmt.Println(res)
		case b == 'd':
			res, newIndex, err := d.Decoders[b](i)

			if err != nil {
				break
			}

			i = newIndex
			fmt.Println(res)
		case b >= '0' && b <= '9':
			res, newIndex, err := d.Decoders[b](i)

			if err != nil {
				break
			}

			i = newIndex

			fmt.Println(res)
		}

		//i++
	}
	return parseDictionaryData(mainMap), nil
}

func parseDictionaryData(mainMap map[any]any) *TorrentFile {
	torrentData := &TorrentFile{}

	if info, exists := mainMap["announce"]; exists {
		torrentData.Announce = info.(string)
	}
	if announce_list, exists := mainMap["announce-list"]; exists {
		torrentData.AnnounceList = announce_list.([][]string)
	}
	if comment, exists := mainMap["comment"]; exists {
		torrentData.Comment = comment.(string)
	}
	if created_by, exists := mainMap["created by"]; exists {
		torrentData.CreatedBy = created_by.(string)
	}
	if creation_date, exists := mainMap["creation date"]; exists {
		torrentData.CreationDate = creation_date.(int)
	}

	if infoDict, exists := mainMap["info"]; exists {
		inner := infoDict.(map[any]any)
		if name, ok := inner["name"]; ok {
			torrentData.Info.Name = name.(string)
		}
		if length, ok := inner["length"]; ok {
			torrentData.Info.Length = length.(int)
		}
		if pieceLength, ok := inner["piece length"]; ok {
			torrentData.Info.PieceLength = pieceLength.(int)
		}
		//pieces
	}

	return torrentData
}

func (d *Decoder) DecodeInt(index int) (any, int, error) {
	end := index
	for i := index; i < len(d.Buffer); i++ {
		b := (d.Buffer)[i]
		if b == 'e' {
			end = i
			break
		}
	}
	num, err := strconv.Atoi(string(d.Buffer[index+1 : end]))
	return num, end + 1, err
}

func (d *Decoder) DecodeString(index int) (any, int, error) {
	end := index
	for i := index; i < len(d.Buffer); i++ {
		b := (d.Buffer)[i]
		if b == ':' {
			end = i
			break
		}
	}
	num, err := strconv.Atoi(string(d.Buffer[index:end]))

	if err != nil {
		return "", -1, err
	}

	return string(d.Buffer[end+1 : end+1+num]), end + num + 1, nil
}

func (d *Decoder) DecodeList(index int) (any, int, error) {
	res := []any{}

	end := index
	for i := index + 1; i < len(d.Buffer); {
		b := (d.Buffer)[i]

		item, newIndex, err := d.Decoders[b](i)

		if err != nil {
			return nil, -1, err
		}

		res = append(res, item)
		i = newIndex

		if (d.Buffer)[i] == 'e' {
			end = i
			break
		}
	}
	return res, end + 1, nil
}

func (d *Decoder) DecodeDictionary(index int) (any, int, error) {
	res := map[any]any{}

	end := index
	isKey := true
	var lastKey any
	for i := index + 1; i < len(d.Buffer); {
		b := (d.Buffer)[i]
		item, newIndex, err := d.Decoders[b](i)

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

		if (d.Buffer)[i] == 'e' {
			end = i
			break
		}
	}

	return res, end + 1, nil
}

func (d *Decoder) DecodeEnd(index int) (any, int, error) {
	return nil, index, nil
}
