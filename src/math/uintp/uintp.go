package uintp

import (
	"math/bits"
)

type UintP struct {
	P     uint64
	value []uint64
}

func New(p uint64) *UintP {
	if p%64 != 0 {
		panic("p must be a multiple of 64")
	}

	return &UintP{
		P:     p,
		value: make([]uint64, p/64),
	}
}

func FromUint(p uint64, u uint64) *UintP {
	r := New(p)
	r.value[0] = u

	return r
}

func FromBytes(p uint64, bs []byte) *UintP {
	r := New(p)
	for i := range r.value {
		r.value[i] = uint64(bs[i*8+0]) |
			uint64(bs[i*8+1])<<8 |
			uint64(bs[i*8+2])<<16 |
			uint64(bs[i*8+3])<<24 |
			uint64(bs[i*8+4])<<32 |
			uint64(bs[i*8+5])<<40 |
			uint64(bs[i*8+6])<<48 |
			uint64(bs[i*8+7])<<56
	}

	return r

}

func Clone(u *UintP) *UintP {
	return &UintP{
		P:     u.P,
		value: append([]uint64{}, u.value...),
	}
}

func (u *UintP) Add(v *UintP) *UintP {
	carry := uint64(0)
	for i := range u.value {
		u.value[i], carry = bits.Add64(u.value[i], v.value[i], carry)
	}

	return u
}

func (u *UintP) AddBytes(bs []byte) *UintP {
	carry := uint64(0)

	for i := range u.value {
		toAdd := uint64(bs[i*8+0]) |
			uint64(bs[i*8+1])<<8 |
			uint64(bs[i*8+2])<<16 |
			uint64(bs[i*8+3])<<24 |
			uint64(bs[i*8+4])<<32 |
			uint64(bs[i*8+5])<<40 |
			uint64(bs[i*8+6])<<48 |
			uint64(bs[i*8+7])<<56
		u.value[i], carry = bits.Add64(u.value[i], toAdd, carry)
	}

	return u
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
	r := New(u.P)
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
	s := shift

	carry := uint64(0)
	carryNew := uint64(0)
	for i := range u.value {
		if s >= 64 {
			u.value[i] = 0
			s -= 64
			continue
		}

		carryNew = u.value[i] >> (64 - s)
		u.value[i] = (u.value[i] << s) | carry
		carry = carryNew
	}

	return u
}

func (u *UintP) Bytes() []byte {
	r := make([]byte, len(u.value)*8)
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
