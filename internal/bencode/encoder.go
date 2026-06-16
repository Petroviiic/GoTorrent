package bencode

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"sort"
	"strconv"
)

// 	Announce     string     `bencode:"announce"`
// 	AnnounceList [][]string `bencode:"announce-list"`
// 	Comment      string     `bencode:"comment"`
// 	CreatedBy    string     `bencode:"created by"`
// 	CreationDate int        `bencode:"creation date"`
// 	Info         InfoDict   `bencode:"info"`

// type InfoDict struct {
// 	Length      int    `bencode:"length"`
// 	Name        string `bencode:"name"`
// 	PieceLength int    `bencode:"piece length"`
// 	Pieces      []byte `bencode:"pieces"`
// }

func Encode(data any) ([]byte, error) {
	var buffer bytes.Buffer

	switch v := data.(type) {
	case int:
		buffer.WriteByte('i')
		buffer.WriteString(strconv.Itoa(v))
		buffer.WriteByte('e')

	case string:
		buffer.WriteString(strconv.Itoa(len(v)))
		buffer.WriteByte(':')
		buffer.WriteString(v)
	case []byte:
		buffer.WriteString(strconv.Itoa(len(v)))
		buffer.WriteByte(':')
		buffer.Write(v)
	case []any:
		buffer.WriteByte('l')
		for _, val := range v {
			data, err := Encode(val)

			if err != nil {
				return nil, err
			}

			buffer.Write(data)
		}
		buffer.WriteByte('e')

	case map[any]any:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k.(string))
		}
		sort.Strings(keys)

		buffer.WriteByte('d')
		for _, key := range keys {
			data, err := Encode(key)
			if err != nil {
				return nil, err
			}
			buffer.Write(data)

			data, err = Encode(v[key])
			if err != nil {
				return nil, err
			}
			buffer.Write(data)
		}
		buffer.WriteByte('e')

	default:
		return nil, fmt.Errorf("unknown data type; Type: %T\n", v)
	}
	return buffer.Bytes(), nil
}

func Hash(data []byte) []byte {
	hash := sha1.New()
	hash.Write(data)
	hashedData := hash.Sum(nil)

	return hashedData
}
