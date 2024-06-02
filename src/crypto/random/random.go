package random

import (
	"crypto/rand"
	"drm-blockchain/src/math/uintp"
	"io"
)

func GenerateBytes(lengthBytes int) ([]byte, error) {
	bs := make([]byte, lengthBytes)
	_, err := io.ReadFull(rand.Reader, bs)
	return bs, err
}

func GenerateUintp(modulusBitsize uint64) (*uintp.UintP, error) {
	lengthBytes := int(modulusBitsize / 8)
	bs, err := GenerateBytes(lengthBytes)
	if err != nil {
		return nil, err
	}

	return uintp.FromBytes(modulusBitsize, bs), nil
}
