package lthash

import (
	"drm-blockchain/src/collections/structures/list"
	"drm-blockchain/src/crypto/hash"
	"drm-blockchain/src/math/algebra/group"
	"drm-blockchain/src/math/algebra/group/bigintmodg"
	"drm-blockchain/src/math/algebra/group/kdirectsumg"
	"math/big"

	"io"

	"golang.org/x/crypto/blake2b"
)

var _ hash.HomHashAlgorithm = &LtHash{}

type opGroup group.Elem[kdirectsumg.KDirectSumG[bigintmodg.BigIntModG]]

func newOpGroup(k int, mod *big.Int) opGroup {
	z := bigintmodg.FromBigInt(big.NewInt(0), mod)
	return opGroup(kdirectsumg.New[bigintmodg.BigIntModG](k, func(i int) group.Elem[bigintmodg.BigIntModG] { return z }))
}

type LtHash struct {
	state      opGroup
	mod        *big.Int
	k          int
	p_bits     int
	block_size int
	xof        blake2b.XOF
}

func buildMask(p_bits int) *big.Int {
	return (&big.Int{}).Lsh(
		big.NewInt(1),
		uint(p_bits),
	)
}

func New(k int, p_bits int, block_size int, key []byte) LtHash {
	xof, err := blake2b.NewXOF(uint32(k*p_bits), key)
	if err != nil {
		panic(err)
	}
	mask := buildMask(p_bits)

	return LtHash{
		state:      newOpGroup(k, mask),
		mod:        mask,
		k:          k,
		p_bits:     p_bits,
		block_size: block_size,
		xof:        xof,
	}
}

func (hash LtHash) randomize(bytes []byte) opGroup {
	hash.xof.Reset()
	hash.xof.Write(bytes)

	buf := make([]byte, hash.k*hash.p_bits)
	_, err := io.ReadFull(hash.xof, buf)
	if err != nil {
		panic(err)
	}

	z := bigintmodg.FromBigInt(big.NewInt(0), hash.mod)
	offset := 0
	entries := kdirectsumg.New[bigintmodg.BigIntModG](hash.k, func(i int) group.Elem[bigintmodg.BigIntModG] { return z })
	entries_iter := entries.Entries.GetIterator()

	for entries_iter.HasNext() {
		part := buf[offset:min(offset+hash.p_bits/8, len(buf)-offset)]
		*(entries_iter.GetNext()) = bigintmodg.FromBytes(part, hash.mod)
		offset += len(part)
	}

	return opGroup(entries)
}

func combine(a1 opGroup, a2 opGroup) opGroup {
	return a1.CombineWith(a2)
}

func (hash *LtHash) Add(bytes []byte) {
	hash.state = combine(hash.state, hash.randomize(bytes))
}

func (hash *LtHash) Remove(bytes []byte) {
	hash.state = combine(hash.state, hash.randomize(bytes).Invert())
}

func (hash *LtHash) ComputeDigest(bytes []byte) {
	offset := 0

	l := list.NewFrom(bytes)

	for offset < len(bytes) {
		part := l.Skip(offset).Take(hash.block_size).ToArray()
		hash.Add(part)
		offset += len(part)
	}
}

func (hash LtHash) GetDigest() []byte {
	r := make([]byte, 0)

	entries_iter := hash.state.AsPure().Entries.GetIterator()
	for entries_iter.HasNext() {
		entry := entries_iter.GetNext()
		val := (*entry)
		r = append(r, val.AsPure().Value.Bytes()...)
	}

	return r
}
