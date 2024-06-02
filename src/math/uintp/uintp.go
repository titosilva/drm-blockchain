package uintp

import (
	"encoding/hex"
	"math/bits"
)

// UintP is a big integer with a modulus of 2^ModulusBitsize
// It uses uint64 operations to perform arithmetic operations more efficiently
// Because of this, the log of the modulus must be a multiple of 64
type UintP struct {
	ModulusBitsize uint64
	value          []uint64
}

func New(modBitsize uint64) *UintP {
	if modBitsize%64 != 0 {
		panic("p must be a multiple of 64")
	}

	return &UintP{
		ModulusBitsize: modBitsize,
		value:          make([]uint64, modBitsize/64),
	}
}

func FromUint(p uint64, u uint64) *UintP {
	r := New(p)
	r.value[0] = u

	return r
}

func FromHex(p uint64, s string) *UintP {
	bs, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	// reverse bs
	for i := 0; i < len(bs)/2; i++ {
		bs[i], bs[len(bs)-1-i] = bs[len(bs)-1-i], bs[i]
	}

	return FromBytes(p, bs)
}

func FromBytes(p uint64, bs []byte) *UintP {
	r := New(p)

	for i := range bs {
		if i >= len(r.value)*8 {
			break
		}

		r.value[i/8] |= uint64(bs[i]) << uint64((i%8)*8)
	}

	return r
}

func Clone(u *UintP) *UintP {
	return &UintP{
		ModulusBitsize: u.ModulusBitsize,
		value:          append([]uint64{}, u.value...),
	}
}

func (u *UintP) Add(v *UintP) *UintP {
	carry := uint64(0)
	for i := range u.value {
		u.value[i], carry = bits.Add64(u.value[i], v.value[i], carry)
	}

	return u
}

func (u *UintP) Mul(v *UintP) *UintP {
	f := FromUint(u.ModulusBitsize, 0)

	for i := range v.value {
		r := Clone(u)
		r.MulUint(v.value[i])
		r.ShiftLeft(uint64(i) * 64)
		f.Add(r)
	}

	u.value = f.value
	return u
}

func (u *UintP) MulUint(v uint64) *UintP {
	carry := uint64(0)

	for i := range u.value {
		nxt_carry, lo := bits.Mul64(u.value[i], v)
		u.value[i] = lo + carry
		carry = nxt_carry
	}

	return u
}

func (u *UintP) AddBytes(bs []byte) *UintP {
	return u.Add(FromBytes(u.ModulusBitsize, bs[:]))
}

func (u *UintP) AddUint(v uint64) *UintP {
	var carry uint64
	u.value[0], carry = bits.Add64(u.value[0], v, 0)
	for i := range u.value[1:] {
		u.value[i+1], carry = bits.Add64(u.value[i+1], 0, carry)
	}

	return u
}

func (u *UintP) Sub(v *UintP) *UintP {
	borrow := uint64(0)
	for i := range u.value {
		u.value[i], borrow = bits.Sub64(u.value[i], v.value[i], borrow)
	}

	return u
}

func (u *UintP) SubBytes(bs []byte) *UintP {
	borrow := uint64(0)

	for i := range u.value {
		toSub := uint64(bs[i*8+0]) |
			uint64(bs[i*8+1])<<8 |
			uint64(bs[i*8+2])<<16 |
			uint64(bs[i*8+3])<<24 |
			uint64(bs[i*8+4])<<32 |
			uint64(bs[i*8+5])<<40 |
			uint64(bs[i*8+6])<<48 |
			uint64(bs[i*8+7])<<56
		u.value[i], borrow = bits.Sub64(u.value[i], toSub, borrow)
	}

	return u
}

func (u *UintP) Inverse() *UintP {
	r := New(u.ModulusBitsize)
	for i := range u.value {
		r.value[i] = ^u.value[i]
	}

	r.AddUint(1)
	return r
}

func (u *UintP) Equals(v *UintP) bool {
	for i := range u.value {
		if u.value[i] != v.value[i] {
			return false
		}
	}

	return true
}

func (u *UintP) ShiftLeft(shift uint64) *UintP {
	if shift == 0 {
		return u
	}

	for i := len(u.value) - 1; i >= 0; i-- {
		if i >= int(shift/64) {
			u.value[i] = u.value[i-int(shift/64)] << (shift % 64)
		} else {
			u.value[i] = 0
		}

		if i > int(shift/64) {
			u.value[i] |= u.value[i-int(shift/64)-1] >> (64 - (shift % 64))
		}
	}

	return u
}

func (u *UintP) Bytes() []byte {
	r := make([]byte, u.ModulusBitsize/8)

	for i, v := range u.value {
		r[i*8+0] = byte(v >> 0)
		r[i*8+1] = byte(v >> 8)
		r[i*8+2] = byte(v >> 16)
		r[i*8+3] = byte(v >> 24)
		r[i*8+4] = byte(v >> 32)
		r[i*8+5] = byte(v >> 40)
		r[i*8+6] = byte(v >> 48)
		r[i*8+7] = byte(v >> 56)
	}

	return r
}

func (u *UintP) Bit(index uint64) byte {
	if index >= u.ModulusBitsize {
		panic("index out of range")
	}

	return byte(u.value[index/64] & (1 << (index % 64)))
}

func (u *UintP) SetBit(index uint64, bit bool) *UintP {
	if index >= u.ModulusBitsize {
		panic("index out of range")
	}

	if bit {
		u.value[index/64] |= 1 << (index % 64)
	} else {
		u.value[index/64] &= ^(1 << (index % 64))
	}

	return u
}
