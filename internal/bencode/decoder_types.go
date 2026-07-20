package bencode

type Decoder struct {
	Buffer   []byte
	Decoders map[byte]func(int) (any, int, error)
}

func NewDecoder(buffer []byte) *Decoder {
	decoder := &Decoder{
		Buffer:   buffer,
		Decoders: make(map[byte]func(int) (any, int, error)),
	}

	decoder.Decoders['i'] = decoder.DecodeInt
	decoder.Decoders['l'] = decoder.DecodeList
	decoder.Decoders['d'] = decoder.DecodeDictionary
	decoder.Decoders['e'] = decoder.DecodeEnd
	for c := '0'; c <= '9'; c++ {
		decoder.Decoders[byte(c)] = decoder.DecodeString
	}

	return decoder
}
