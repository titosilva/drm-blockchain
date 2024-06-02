package gcrypto

import (
	"drm-blockchain/src/crypto/random"
	"drm-blockchain/src/crypto/random/drbg/sha256drbg"
	"drm-blockchain/src/math/uintp"
)

type GCrypto struct {
	modulusBitsize uint64
}

func New(modulusBitsize uint64) *GCrypto {
	r := new(GCrypto)

	r.modulusBitsize = modulusBitsize

	return r
}

func (g *GCrypto) Encrypt(data []byte, key []byte) []byte {
	encoded := g.Encode(data)
	encrypted := g.EncryptEncoded(encoded, key)

	return toBytes(encrypted)
}

func (g *GCrypto) Decrypt(data []byte, key []byte) []byte {
	decryptedEncoded := g.DecryptEncoded(fromBytes(data, g.modulusBitsize), key)
	decoded := g.Decode(decryptedEncoded)

	return decoded
}

func (g *GCrypto) Encode(data []byte) []*uintp.UintP {
	// Encodes each bit in data to a randomly generated number,
	// being even or odd depending on the bit value
	r := make([]*uintp.UintP, len(data)*8)

	for i := 0; i < len(data)*8; i++ {
		u, err := random.GenerateUintp(g.modulusBitsize)
		if err != nil {
			panic(err)
		}

		u.SetBit(0, data[i/8]&(1<<(i%8)) != 0)
		r[i] = u
	}

	return r
}

func (g *GCrypto) Decode(encodedData []*uintp.UintP) []byte {
	// Decodes each number in encodedData to a bit
	r := make([]byte, (len(encodedData)+7)/8)

	for i := range encodedData {
		r[i/8] |= byte(encodedData[i].Bit(0)) << uint(i%8)
	}

	return r
}

func (g *GCrypto) EncryptEncoded(encodedData []*uintp.UintP, key []byte) []*uintp.UintP {
	r := make([]*uintp.UintP, len(encodedData))

	drbg := sha256drbg.New()
	drbg.Seed(key)

	for i := range encodedData {
		generated, _ := drbg.Generate(int(g.modulusBitsize / 8))
		r[i] = uintp.Clone(encodedData[i])
		r[i].Add(uintp.FromBytes(g.modulusBitsize, generated))
	}

	return r
}

func (g *GCrypto) DecryptEncoded(encryptedData []*uintp.UintP, key []byte) []*uintp.UintP {
	r := make([]*uintp.UintP, len(encryptedData))

	drbg := sha256drbg.New()
	drbg.Seed(key)

	for i := range encryptedData {
		generated, _ := drbg.Generate(int(g.modulusBitsize / 8))
		r[i] = uintp.Clone(encryptedData[i])
		r[i].Sub(uintp.FromBytes(g.modulusBitsize, generated))
	}

	return r
}

func toBytes(data []*uintp.UintP) []byte {
	r := make([]byte, (len(data)*int(data[0].ModulusBitsize)+7)/8)

	for i := range data {
		bs := data[i].Bytes()
		for j := range bs {
			r[i*len(bs)+j] = bs[j]
		}
	}

	return r
}

func fromBytes(data []byte, modulusBitsize uint64) []*uintp.UintP {
	r := make([]*uintp.UintP, (len(data)*8+int(modulusBitsize)-1)/int(modulusBitsize))

	for i := 0; i < len(r); i++ {
		r[i] = uintp.FromBytes(modulusBitsize, data[i*int(modulusBitsize/8):])
	}

	return r
}
