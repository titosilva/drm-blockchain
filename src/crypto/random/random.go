package random

import (
	"crypto/rand"
	"io"
)

func GenerateBytes(length int) ([]byte, error) {
	bs := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, bs)
	return bs, err
}
