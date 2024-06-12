package gcrypt

import (
	"crypto/sha256"
	"drm-blockchain/src/crypto/random/drbg/sha256drbg"
	"drm-blockchain/src/math/uintp"
)

type GCrypt struct {
	modulusBitsize uint64
}

func New(modulusBitsize uint64) *GCrypt {
	r := new(GCrypt)

	r.modulusBitsize = modulusBitsize

	return r
}

func (g *GCrypt) Encrypt(data []byte, key []byte) []byte {
	encoded := g.Encode(data)
	encrypted := g.EncryptEncoded(encoded, key)

	return g.ToBytes(encrypted)
}

func (g *GCrypt) Decrypt(data []byte, key []byte) []byte {
	decryptedEncoded := g.DecryptEncoded(g.FromBytes(data), key)
	decoded := g.Decode(decryptedEncoded)

	return decoded
}

func (g *GCrypt) EncodeToBytes(data []byte) []byte {
	encoded := g.Encode(data)

	return g.ToBytes(encoded)
}

func (g *GCrypt) Encode(data []byte) []*uintp.UintP {
	// Encodes each bit in data to a randomly generated number,
	// being even or odd depending on the bit value
	r := make([]*uintp.UintP, len(data)*8)
	drbg := sha256drbg.New()
	seed := sha256.Sum256(data)
	drbg.Seed(seed[:])

	for i := 0; i < len(data)*8; i++ {
		bs, err := drbg.Generate(int(g.modulusBitsize) / 8)
		if err != nil {
			panic(err)
		}

		u := uintp.FromBytes(g.modulusBitsize, bs)
		u.SetBit(0, data[i/8]&(1<<(i%8)) != 0)
		r[i] = u
	}

	return r
}

func (g *GCrypt) ExpandKey(key []byte, lengthBlocks int) []*uintp.UintP {
	// Expands the key to the desired length using a DRBG
	drbg := sha256drbg.New()
	drbg.Seed(key)

	r := make([]*uintp.UintP, lengthBlocks)

	for i := 0; i < lengthBlocks; i++ {
		generated, _ := drbg.Generate(int(g.modulusBitsize / 8))
		r[i] = uintp.FromBytes(g.modulusBitsize, generated)
	}

	return r
}

func (g *GCrypt) ExpandKeyToBytes(key []byte, lengthBlocks int) []byte {
	expandedKey := g.ExpandKey(key, lengthBlocks)

	return g.ToBytes(expandedKey)
}

func (g *GCrypt) Decode(encodedData []*uintp.UintP) []byte {
	// Decodes each number in encodedData to a bit
	r := make([]byte, (len(encodedData)+7)/8)

	for i := range encodedData {
		r[i/8] |= byte(encodedData[i].Bit(0)) << uint(i%8)
	}

	return r
}

func (g *GCrypt) EncryptEncoded(encodedData []*uintp.UintP, key []byte) []*uintp.UintP {
	r := make([]*uintp.UintP, len(encodedData))
	expandedKey := g.ExpandKey(key, len(encodedData))

	for i := range encodedData {
		r[i] = uintp.Clone(encodedData[i])
		r[i].Add(expandedKey[i])
	}

	return r
}

func (g *GCrypt) DecryptEncoded(encryptedData []*uintp.UintP, key []byte) []*uintp.UintP {
	r := make([]*uintp.UintP, len(encryptedData))
	expandedKey := g.ExpandKey(key, len(encryptedData))

	for i := range encryptedData {
		r[i] = uintp.Clone(encryptedData[i])
		r[i].Sub(expandedKey[i])
	}

	return r
}

func (g GCrypt) ToBytes(data []*uintp.UintP) []byte {
	r := make([]byte, len(data)*int(g.modulusBitsize)/8)

	for i := range data {
		bs := data[i].Bytes()
		for j := range bs {
			r[i*len(bs)+j] = bs[j]
		}
	}

	return r
}

func (g GCrypt) FromBytes(data []byte) []*uintp.UintP {
	r := make([]*uintp.UintP, (len(data)*8+int(g.modulusBitsize)-1)/int(g.modulusBitsize))

	for i := 0; i < len(r); i++ {
		r[i] = uintp.FromBytes(g.modulusBitsize, data[i*int(g.modulusBitsize/8):])
	}

	return r
}
