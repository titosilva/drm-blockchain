package uintp_test

import (
	"drm-blockchain/src/math/uintp"
	ez "drm-blockchain/src/utils/test"
	"testing"
)

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

func Test__Uintp__MulUint__PowerOf2WithOverflow__ShouldEqual__HardcodedResult(t *testing.T) {
	u := uintp.FromBytes(128, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	v := uintp.FromBytes(128, []byte{0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01})

	r := u.MulUint(uint64(2))

	ez := ez.New(t)
	ez.Assert(r.Equals(v))
}

func Test__Uintp__Mul__PowerOf2__ShouldEqual__ShiftLeft(t *testing.T) {
	u := uintp.FromUint(128, 1)
	v := uintp.FromUint(128, 1<<32)

	u_cp := uintp.Clone(u)

	ez := ez.New(t)
	ez.Assert(u.Mul(v).Equals(u_cp.ShiftLeft(32)))
}

func Test__Uintp__Mul__LargeNumbers__ShouldEqual__RightAnswer(t *testing.T) {
	ez := ez.New(t)

	u := uintp.FromBytes(128, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	v := uintp.FromBytes(128, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	exp := uintp.FromHex(128, "fffffffffffffffe0000000000000001")

	r := u.Mul(v)
	ez.Assert(r.Equals(exp))
}
