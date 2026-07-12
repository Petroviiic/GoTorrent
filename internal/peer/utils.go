package peer

import "fmt"

func (p *PeerClient) HasPiece(pieceIndex int) bool {
	if p.Bitfield == nil {
		return false
	}

	//bitfield example: 255 255 255 255 0 = 11111111 11111111 11111111 11111111 00000000
	byteIndex := pieceIndex / 8

	if byteIndex >= len(p.Bitfield) {
		return false
	}

	bit := pieceIndex % 8
	bitmask := 128 //10000000

	for i := 0; i < bit; i++ {
		bitmask >>= 1
	}

	return p.Bitfield[byteIndex]&byte(bitmask) != 0
}

func (p *PeerClient) UpdatePiece(pieceIndex int) {
	if p.Bitfield == nil {
		totalPieces := p.Manager.TotalPieces
		p.Bitfield = make([]byte, (totalPieces+7)/8)
	}

	//bitfield example: 255 255 255 255 0 = 11111111 11111111 11111111 11111111 00000000
	byteIndex := pieceIndex / 8

	if byteIndex >= len(p.Bitfield) {
		return
	}

	bit := pieceIndex % 8
	bitmask := 128 //10000000

	for i := 0; i < bit; i++ {
		bitmask >>= 1
	}

	p.Bitfield[byteIndex] |= byte(bitmask)
}

func HashOk(downloadedPieces []*PieceOfResult, expected []byte) ([]byte, bool) {
	//gotHash := []byte{}
	fmt.Println("evo me", downloadedPieces, expected)
	for _, piece := range downloadedPieces {
		fmt.Println(len(piece.Downloaded))
	}

	return nil, false
}
