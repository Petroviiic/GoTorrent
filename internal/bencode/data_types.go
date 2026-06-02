package bencode

type Decoder struct {
	buffer   []byte
	decoders map[byte]func(int) (any, int, error) // Funkcije sad primaju samo indeks!
}

func newDecoder(buffer []byte) *Decoder {
	decoder := &Decoder{
		buffer: buffer,
	}

	decoder.decoders['i'] = decoder.decodeInt
	decoder.decoders['l'] = decoder.decodeList
	// DecoderFunc['d'] = decodeDictionary
	for c := '0'; c <= '9'; c++ {
		decoder.decoders[byte(c)] = decoder.decodeString
	}

	return decoder
}
