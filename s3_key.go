package main

import (
	"crypto/rand"
	"encoding/hex"
)

func makeKey() (string, error) {
	byteSlice := make([]byte, 32)
	rand.Read(byteSlice)
	key := hex.EncodeToString(byteSlice)

	return key + ".mp4", nil
}
