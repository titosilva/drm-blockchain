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
	decryptedEncoded := g.DecryptEncoded(g.fromBytes(data), key)
	decoded := g.Decode(decryptedEncoded)

	return decoded
}

func (g *GCrypto) EncodeToBytes(data []byte) []byte {
	encoded := g.Encode(data)

	return toBytes(encoded)
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

func (g *GCrypto) ExpandKey(key []byte, length int) []*uintp.UintP {
	// Expands the key to the desired length using a DRBG
	drbg := sha256drbg.New()
	drbg.Seed(key)

	r := make([]*uintp.UintP, length)

	for i := 0; i < length; i++ {
		generated, _ := drbg.Generate(int(g.modulusBitsize / 8))
		r[i] = uintp.FromBytes(g.modulusBitsize, generated)
	}

	return r
}

func (g *GCrypto) EncodeKeyToBytes(key []byte, length int) []byte {
	expandedKey := g.ExpandKey(key, length)

	return toBytes(expandedKey)
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
	expandedKey := g.ExpandKey(key, len(encodedData))

	for i := range encodedData {
		r[i] = uintp.Clone(encodedData[i])
		r[i].Add(expandedKey[i])
	}

	return r
}

func (g *GCrypto) DecryptEncoded(encryptedData []*uintp.UintP, key []byte) []*uintp.UintP {
	r := make([]*uintp.UintP, len(encryptedData))
	expandedKey := g.ExpandKey(key, len(encryptedData))

	for i := range encryptedData {
		r[i] = uintp.Clone(encryptedData[i])
		r[i].Sub(expandedKey[i])
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

func (g GCrypto) fromBytes(data []byte) []*uintp.UintP {
	r := make([]*uintp.UintP, (len(data)*8+int(g.modulusBitsize)-1)/int(g.modulusBitsize))

	for i := 0; i < len(r); i++ {
		r[i] = uintp.FromBytes(g.modulusBitsize, data[i*int(g.modulusBitsize/8):])
	}

	return r
}
