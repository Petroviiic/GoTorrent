package bencode

type Decoders map[byte]func([]byte, int) (any, int, error)

func loadDecoders() *Decoders {
	decoders := make(Decoders)
	decoders['i'] = decodeInt
	decoders['l'] = decodeList
	// DecoderFunc['d'] = decodeDictionary
	for c := '0'; c <= '9'; c++ {
		decoders[byte(c)] = decodeString
	}

	return &decoders
}
