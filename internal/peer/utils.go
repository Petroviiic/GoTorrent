package peer

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
)

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
	gotData := []byte{}
	for _, piece := range downloadedPieces {
		gotData = append(gotData, piece.Downloaded...)
	}

	sum := sha1.Sum(gotData)
	if bytes.Equal(sum[:], expected) {
		return gotData, true
	}

	return nil, false
}

func DecodePiece(data []byte) *PieceOfResult {
	index := binary.BigEndian.Uint32(data[0:4])
	begin := binary.BigEndian.Uint32(data[4:8])
	//length := binary.BigEndian.Uint32(data[8:12])	//not sure bout this one

	fmt.Println(index, begin, len(data[8:]))
	return &PieceOfResult{
		PieceIndex:  int(index),
		BlockOffset: int(begin),
		Downloaded:  data[8:],
	}
}
