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
	lthash     *lthash.LtHash
	nonceHash  []byte
	nonceState []*uintp.UintP
	key        []byte
}

func New(chunk_count uint, chunk_size_bits uint, block_size_bytes int, key []byte) *GHash {
	r := new(GHash)
	r.lthash = lthash.New(chunk_count, chunk_size_bits, block_size_bytes, key)
	r.key = key

	return r
}

func (hash *GHash) SetNonce(nonce []byte) {
	hash.lthash.Reset()
	hash.lthash.Add(nonce)
	hash.nonceHash = hash.lthash.GetDigest()
	hash.nonceState = hash.lthash.GetState()
}

func (hash *GHash) SetNonceHash(nonceHash []byte) {
	hash.lthash.Reset()
	hash.lthash.CombineBytes(nonceHash)
	hash.nonceHash = nonceHash
	hash.nonceState = hash.lthash.GetState()
}

func (hash *GHash) RemoveNonce(nonce []byte) {
	hash.lthash.Remove(nonce)
}

func (hash *GHash) RemoveNonceState(nonceState []*uintp.UintP) {
	hash.lthash.CombineInverse(nonceState)
}

func (hash *GHash) SetNonceState(nonceState []*uintp.UintP) {
	hash.lthash.Reset()
	hash.lthash.Combine(nonceState)
	hash.nonceHash = hash.lthash.GetDigest()
	hash.nonceState = nonceState
}

func (hash *GHash) GetNonceHash() []byte {
	r := make([]byte, len(hash.nonceHash))

	for i := range hash.nonceHash {
		r[i] = hash.nonceHash[i]
	}

	return r
}

func (hash *GHash) GetNonceState() []*uintp.UintP {
	r := make([]*uintp.UintP, len(hash.nonceState))

	for i := range hash.nonceState {
		r[i] = uintp.Clone(hash.nonceState[i])
	}

	return r
}

func (hash *GHash) GetState() []*uintp.UintP {
	return hash.lthash.GetState()
}

func (hash *GHash) AddBytes(data []byte) {
	for i := 0; i < len(data); i += int(hash.lthash.ModulusBitsize / 8) {
		block := data[i:min(i+int(hash.lthash.ModulusBitsize/8), len(data))]
		hash.AddBlockWithIndex(block, uint(i))
	}
}

func (hash *GHash) RemoveBytes(data []byte) {
	for i := 0; i < len(data); i += int(hash.lthash.ModulusBitsize / 8) {
		block := data[i:min(i+int(hash.lthash.ModulusBitsize/8), len(data))]
		hash.RemoveBlockWithIndex(block, uint(i))
	}
}

func (hash *GHash) Add(u []*uintp.UintP) {
	hash.lthash.Combine(u)
}

func (hash *GHash) Remove(u []*uintp.UintP) {
	hash.lthash.CombineInverse(u)
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
