package utils

import (
	"crypto/rand"
)

func GeneratePeerID(fixedKey []byte) []byte {
	peerID := make([]byte, 20)

	rand.Read(peerID)

	copy(peerID[:len(fixedKey)], fixedKey)
	return peerID
}
