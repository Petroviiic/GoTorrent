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
	Length      int               `bencode:"length"`
	Name        string            `bencode:"name"`
	PieceLength int               `bencode:"piece length"`
	Pieces      []byte            `bencode:"pieces"`
	Files       []TorrentFileItem `bencode:"files"`
}

type TorrentFileItem struct {
	Length int
	Path   []string
}

func LoadAndDecode(path string) (*TorrentFile, []byte, error) {
	//cekiraj jel validan path

	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	decoder := NewDecoder(buffer)

	torrentDataMap, err := decoder.Decode(buffer, 0)
	if err != nil {
		return nil, nil, err
	}

	if _, exists := torrentDataMap["info"]; !exists {
		return nil, nil, fmt.Errorf("info field doesnt exist")
	}

	encodedInfoDict, err := Encode(torrentDataMap["info"])
	if err != nil {
		return nil, nil, err
	}

	return ParseTorrentMap(torrentDataMap), Hash(encodedInfoDict), err
}

func (d *Decoder) Decode(buffer []byte, index int) (map[any]any, error) {
	mainMap := map[any]any{}
	for i := index; i < len(buffer); {
		switch b := buffer[i]; {
		case b == 'd':
			res, newIndex, err := d.Decoders[b](i)

			if err != nil {
				return nil, err
			}

			i = newIndex
			mainMap = res.(map[any]any)
			break
		case b == '\r' || b == '\n' || b == ' ' || b == '\t':
			i++
			break
		}
	}

	return mainMap, nil
}

func ParseTorrentMap(mainMap map[any]any) *TorrentFile {
	torrentData := &TorrentFile{}

	if info, exists := mainMap["announce"]; exists {
		torrentData.Announce = string(info.([]byte))
	}
	if announce_list, exists := mainMap["announce-list"]; exists {

		for _, outer := range announce_list.([]interface{}) {
			items := []string{}
			for _, inner := range outer.([]interface{}) {
				items = append(items, string(inner.([]byte)))
			}

			torrentData.AnnounceList = append(torrentData.AnnounceList, items)
		}
	}
	if comment, exists := mainMap["comment"]; exists {
		torrentData.Comment = string(comment.([]byte))
	}
	if created_by, exists := mainMap["created by"]; exists {
		torrentData.CreatedBy = string(created_by.([]byte))
	}
	if creation_date, exists := mainMap["creation date"]; exists {
		torrentData.CreationDate = creation_date.(int)
	}

	if infoDict, exists := mainMap["info"]; exists {
		inner := infoDict.(map[any]any)
		if name, ok := inner["name"]; ok {
			torrentData.Info.Name = string(name.([]byte))
		}
		if pieceLength, ok := inner["piece length"]; ok {
			torrentData.Info.PieceLength = pieceLength.(int)
		}
		if pieces, ok := inner["pieces"]; ok {
			torrentData.Info.Pieces = pieces.([]byte)
		}

		if files, ok := inner["files"]; ok {
			for _, file := range files.([]interface{}) {
				fileMap := file.(map[any]any)

				torrentFileItem := TorrentFileItem{
					Length: fileMap["length"].(int),
				}

				for _, path := range fileMap["path"].([]interface{}) {
					torrentFileItem.Path = append(torrentFileItem.Path, string(path.([]byte)))
				}
				torrentData.Info.Files = append(torrentData.Info.Files, torrentFileItem)
			}
		} else {
			if length, ok := inner["length"]; ok {
				torrentData.Info.Length = length.(int)
			}
		}
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

	return d.Buffer[end+1 : end+1+num], end + num + 1, nil
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
			item = string(item.([]byte))
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
