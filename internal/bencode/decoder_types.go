package bencode

import (
	"fmt"
	"io"
)

type Decoder struct {
	Buffer   []byte
	decoders map[byte]func(int) (any, int, error)
}

func NewDecoder(buffer []byte) *Decoder {
	decoder := &Decoder{
		Buffer:   buffer,
		decoders: make(map[byte]func(int) (any, int, error)),
	}

	decoder.decoders['i'] = decoder.DecodeInt
	decoder.decoders['l'] = decoder.DecodeList
	decoder.decoders['d'] = decoder.DecodeDictionary
	decoder.decoders['e'] = decoder.DecodeEnd
	for c := '0'; c <= '9'; c++ {
		decoder.decoders[byte(c)] = decoder.DecodeString
	}

	return decoder
}
func (d *Decoder) DecodeByte(idx int) (any, int, error) {
	if idx >= len(d.Buffer) {
		return nil, idx, io.ErrUnexpectedEOF
	}
	b := d.Buffer[idx]

	f, exists := d.decoders[b]
	if !exists {
		return nil, 0, fmt.Errorf("bencode: invalid or unsupported token '%c' (byte %d) at index %d", b, b, idx)
	}
	return f(idx)
}
