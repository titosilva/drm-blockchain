package lthash

import (
	"drm-blockchain/src/collections/structures/list"
	"drm-blockchain/src/math/uintp"

	"golang.org/x/crypto/blake2b"
)

type LtHash struct {
	ModulusBitsize   uint64
	chunks           []*uintp.UintP
	chunk_count      uint
	chunk_size_bits  uint
	block_size_bytes int
	xof              blake2b.XOF
	chunk_buf        []byte
}

func getChunksWithZero(chunk_bits uint, chunk_count uint) []*uintp.UintP {
	chunks := make([]*uintp.UintP, chunk_count)

	for i := range chunks {
		chunks[i] = uintp.FromUint(uint64(chunk_bits), 0)
	}

	return chunks
}

func New(chunk_count uint, chunk_size_bits uint, block_size_bytes int, key []byte) *LtHash {
	xof, err := blake2b.NewXOF(uint32(chunk_count*chunk_size_bits), key)
	if err != nil {
		panic(err)
	}

	r := new(LtHash)
	r.chunks = getChunksWithZero(chunk_size_bits, chunk_count)
	r.chunk_count = chunk_count
	r.chunk_size_bits = chunk_size_bits
	r.block_size_bytes = block_size_bytes
	r.ModulusBitsize = uint64(chunk_size_bits)
	r.xof = xof
	r.chunk_buf = make([]byte, chunk_size_bits/8)

	return r
}

func NewDirect(chunk_count uint, chunk_size_bits uint, block_size_bytes int, key []byte) LtHash {
	xof, err := blake2b.NewXOF(uint32(chunk_count*chunk_size_bits), key)
	if err != nil {
		panic(err)
	}

	return LtHash{
		chunks:           getChunksWithZero(chunk_size_bits, chunk_count),
		chunk_count:      chunk_count,
		chunk_size_bits:  chunk_size_bits,
		block_size_bytes: block_size_bytes,
		ModulusBitsize:   uint64(chunk_size_bits),
		xof:              xof,
		chunk_buf:        make([]byte, chunk_size_bits/8),
	}
}

func (hash *LtHash) Reset() {
	hash.chunks = getChunksWithZero(hash.chunk_size_bits, hash.chunk_count)
}

func (hash LtHash) randomizeThenCombine(bytes []byte) {
	hash.xof.Reset()
	hash.xof.Write(bytes)

	for i := range hash.chunks {
		_, err := hash.xof.Read(hash.chunk_buf)
		if err != nil {
			panic(err)
		}

		to_add := uintp.FromBytes(uint64(hash.chunk_size_bits), hash.chunk_buf)
		hash.chunks[i].Add(to_add)
	}
}

func (hash LtHash) randomizeThenCombineInverse(bytes []byte) {
	hash.xof.Reset()
	hash.xof.Write(bytes)

	for i := range hash.chunks {
		_, err := hash.xof.Read(hash.chunk_buf)
		if err != nil {
			panic(err)
		}

		hash.chunks[i].SubBytes(hash.chunk_buf)
	}
}

func (hash LtHash) randomizeThenCombineMul(mul *uintp.UintP, bytes []byte) {
	hash.xof.Reset()
	hash.xof.Write(bytes)

	for i := range hash.chunks {
		_, err := hash.xof.Read(hash.chunk_buf)
		if err != nil {
			panic(err)
		}

		to_add := uintp.FromBytes(uint64(hash.chunk_size_bits), hash.chunk_buf)
		to_add.Mul(mul)
		hash.chunks[i].Add(to_add)
	}
}

func (hash *LtHash) Add(bytes []byte) {
	hash.randomizeThenCombine(bytes)
}

func (hash *LtHash) AddMul(mul *uintp.UintP, bytes []byte) {
	hash.randomizeThenCombineMul(mul, bytes)
}

func (hash *LtHash) Remove(bytes []byte) {
	hash.randomizeThenCombineInverse(bytes)
}

func (hash *LtHash) RemoveMul(mul *uintp.UintP, bytes []byte) {
	hash.randomizeThenCombineMul(mul.Inverse(), bytes)
}

func (hash *LtHash) ComputeDigest(bytes []byte) {
	offset := 0

	l := list.NewFrom(bytes)

	for offset < len(bytes) {
		part := l.Skip(offset).Take(hash.block_size_bytes).ToArray()
		hash.Add(part)
		offset += len(part)
	}
}

func (hash LtHash) GetDigest() []byte {
	r := make([]byte, 0)

	for _, chunk := range hash.chunks {
		r = append(r, chunk.Bytes()...)
	}

	return r
}

func (hash LtHash) GetState() []*uintp.UintP {
	r := make([]*uintp.UintP, len(hash.chunks))

	for i := range hash.chunks {
		r[i] = uintp.Clone(hash.chunks[i])
	}

	return hash.chunks
}

func (hash *LtHash) Combine(state []*uintp.UintP) {
	for i := range hash.chunks {
		hash.chunks[i].Add(state[i])
	}
}

func (hash *LtHash) CombineBytes(state []byte) {
	toCombine := make([]*uintp.UintP, hash.chunk_count)

	for i := 0; i < len(state); i += int(hash.ModulusBitsize / 8) {
		block := state[i : i+int(hash.ModulusBitsize)]
		toCombine[i/int(hash.ModulusBitsize/8)] = uintp.FromBytes(hash.ModulusBitsize, block)
	}

	hash.Combine(toCombine)
}

func (hash *LtHash) CombineInverse(state []*uintp.UintP) {
	for i := range hash.chunks {
		hash.chunks[i].Sub(state[i])
	}
}
