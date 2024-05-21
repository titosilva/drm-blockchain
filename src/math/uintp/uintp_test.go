package uintp_test

import (
	"drm-blockchain/src/math/uintp"
	ez "drm-blockchain/src/utils/test"
	"encoding/hex"
	"testing"
)

var testCasesAdd = []struct {
	p         uint64
	u, v, exp string
}{
	{128, "ffffffffffffffff", "cafe", "01000000000000cafd"},
	{128,
		"ffffffffffffffffffffffffffffffff",
		"ffffffffffffffffffffffffffffffff",
		"fffffffffffffffffffffffffffffffe",
	},
}

func Test__Uintp__Add__ShouldEqual__Sum(t *testing.T) {
	for _, tc := range testCasesAdd {
		ez := ez.New(t)

		u := uintp.FromHex(tc.p, tc.u)
		v := uintp.FromHex(tc.p, tc.v)

		exp := uintp.FromHex(tc.p, tc.exp)

		r := u.Add(v)
		ez.Assert(r.Equals(exp))
	}
}

var testCasesShiftLeft = []struct {
	p      uint64
	u, exp string
	shift  uint64
}{
	{64, "ffffffffffffffff", "fffffffffffffe00", 9},
	{128,
		"ffffffffffffffffffffffffffffffff",
		"ffffffffffffffff0000000000000000",
		64,
	},
	{128,
		"ffffffffffffffffffffffffffffffff",
		"ffffffffffffff000000000000000000",
		72,
	},
	{128,
		"ffffffffffffffffffffffffffffffff",
		"fffffffffffffe000000000000000000",
		73,
	},
}

func Test__Uintp__ShiftLeft__ShouldEqual__HardcodedResult(t *testing.T) {
	for _, tc := range testCasesShiftLeft {
		ez := ez.New(t)

		u := uintp.FromHex(tc.p, tc.u)
		exp := uintp.FromHex(tc.p, tc.exp)

		r := u.ShiftLeft(tc.shift)
		ez.Assert(r.Equals(exp))
	}
}

func Test__Uintp__MulUint__PowerOf2__ShouldEqual__ShiftLeft(t *testing.T) {
	u := uintp.FromUint(128, 1)
	v := uintp.Clone(u)

	ez := ez.New(t)
	ez.Assert(u.MulUint(uint64(1 << 3)).Equals(v.ShiftLeft(3)))
}

func Test__Uintp__MulUint__PowerOf2WithOverflow__ShouldEqual__ShiftLeft(t *testing.T) {
	u := uintp.FromBytes(128, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	v := uintp.Clone(u)

	r := u.MulUint(uint64(2))

	ez := ez.New(t)
	ez.Assert(r.Equals(v.ShiftLeft(1)))
}

var testCasesMulUintOverflow = []struct {
	p         uint64
	u, v, exp string
}{
	{128, "ffffffffffffffffffffffffffffffff", "02", "fffffffffffffffffffffffffffffffe"},
	{128,
		"ffffffffffffffffffffffffffffffff",
		"ffffffffffffffff",
		"ffffffffffffffff0000000000000001",
	},
}

func Test__Uintp__MulUint__PowerOf2WithOverflow__ShouldEqual__HardcodedResult(t *testing.T) {
	for _, tc := range testCasesMulUintOverflow {
		u := uintp.FromHex(tc.p, tc.u)
		bs, _ := hex.DecodeString(tc.v)
		v := uint64(0)

		for i := range bs {
			v |= uint64(bs[i]) << uint64(i*8)
		}

		r := u.MulUint(v)

		exp := uintp.FromHex(tc.p, tc.exp)
		ez := ez.New(t)
		ez.Assert(r.Equals(exp))
	}
}

func Test__Uintp__Mul__PowerOf2__ShouldEqual__ShiftLeft(t *testing.T) {
	u := uintp.FromUint(128, 1)
	v := uintp.FromUint(128, 1<<32)

	u_cp := uintp.Clone(u)

	ez := ez.New(t)
	ez.Assert(u.Mul(v).Equals(u_cp.ShiftLeft(32)))
}

var testCases = []struct {
	p         uint64
	u, v, exp string
}{
	{128, "ffffffffffffffff", "cafe", "cafdffffffffffff3502"},
	{128, "ffffffffffffffff", "ffffffffffffffff", "fffffffffffffffe0000000000000001"},
	{128,
		"ffffffffffffffffffffffffffffffff",
		"ffffffffffffffffffffffffffffffff",
		"01",
	},
}

func Test__Uintp__Mul__LargeNumbers__ShouldEqual__RightAnswer(t *testing.T) {
	for _, tc := range testCases {
		ez := ez.New(t)

		u := uintp.FromHex(tc.p, tc.u)
		v := uintp.FromHex(tc.p, tc.v)

		exp := uintp.FromHex(tc.p, tc.exp)

		r := u.Mul(v)
		ez.Assert(r.Equals(exp))
	}
}

func Test__Uintp__FromHex__Should__ConvertCorrectly(t *testing.T) {
	ez := ez.New(t)

	u := uintp.FromHex(128, "01fffffffffffffffe")
	exp := uintp.FromBytes(128, []byte{0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01})
	ez.Assert(u.Equals(exp))
}
