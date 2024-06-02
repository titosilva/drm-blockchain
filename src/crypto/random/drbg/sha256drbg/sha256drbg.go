package sha256drbg

import (
	"crypto/sha256"
	"drm-blockchain/src/crypto/random/drbg"
)

type SHA256DRBG struct {
	seed  []byte
	state []byte
}

// static interface check
var _ drbg.DRBG = (*SHA256DRBG)(nil)

func New() *SHA256DRBG {
	r := new(SHA256DRBG)
	return r
}

// GenerateByte implements drbg.DRBG.
func (h *SHA256DRBG) GenerateByte() (byte, error) {
	bs := sha256.Sum256(h.state)
	h.state = bs[:]
	return bs[0], nil
}

func (h *SHA256DRBG) Generate(bytes int) ([]byte, error) {
	bs := make([]byte, bytes)

	for i := 0; i < bytes; i++ {
		bs[i], _ = h.GenerateByte()
	}

	return bs, nil
}

// Seed implements drbg.DRBG.
func (h *SHA256DRBG) Seed(seed []byte) error {
	h.seed = seed
	h.state = seed
	return nil
}
