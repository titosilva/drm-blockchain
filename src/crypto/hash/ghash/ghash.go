package ghash

import (
	"drm-blockchain/src/crypto/hash/lthash"
	"drm-blockchain/src/math/uintp"
)

// GHash takes blocks of data and use them as multipliers for the lthash of their index.
// The size of each block must be the same size as the lthash modulus
// GHash also must receive a nonce to be used as the error, added to make it one way

type GHash struct {
	// Undeyling Lthash algorithm
	lthash    *lthash.LtHash
	nonceHash []byte
	key       []byte
}

func New(chunk_count uint, chunk_size_bits uint, block_size_bytes int, key []byte, nonce []byte) *GHash {
	r := new(GHash)
	r.lthash = lthash.New(chunk_count, chunk_size_bits, block_size_bytes, key)
	r.key = key

	r.lthash.Add(nonce)
	r.nonceHash = r.lthash.GetDigest()

	return r
}

func NewFromNonceHash(chunk_count uint, chunk_size_bits uint, block_size_bytes int, key []byte, nonceHash []byte) *GHash {
	r := new(GHash)
	r.lthash = lthash.New(chunk_count, chunk_size_bits, block_size_bytes, key)

	r.lthash.CombineBytes(nonceHash)

	r.key = key
	r.nonceHash = nonceHash

	return r
}

func (hash GHash) GetNonceHash() []byte {
	return hash.nonceHash
}

func (hash *GHash) Add(data []byte) {
	for i := 0; i < len(data); i += int(hash.lthash.ModulusBitsize / 8) {
		block := data[i : i+int(hash.lthash.ModulusBitsize)]
		hash.AddBlockWithIndex(block, uint(i))
	}
}

func (hash *GHash) Remove(data []byte) {
	for i := 0; i < len(data); i += int(hash.lthash.ModulusBitsize / 8) {
		block := data[i : i+int(hash.lthash.ModulusBitsize)]
		hash.RemoveBlockWithIndex(block, uint(i))
	}
}

func (hash *GHash) AddBlockWithIndex(block []byte, index uint) {
	mul := uintp.FromBytes(hash.lthash.ModulusBitsize, block)
	hash.lthash.AddMul(mul, []byte{byte(index)})
}

func (hash *GHash) RemoveBlockWithIndex(block []byte, index uint) {
	mul := uintp.FromBytes(hash.lthash.ModulusBitsize, block)
	hash.lthash.RemoveMul(mul, []byte{byte(index)})
}

func (hash GHash) GetDigest() []byte {
	return hash.lthash.GetDigest()
}
