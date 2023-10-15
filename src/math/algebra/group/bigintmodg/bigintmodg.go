package bigintmodg

import (
	"drm-blockchain/src/math/algebra/group"
	"math/big"
)

// static interface check
var _ group.Elem[BigIntModG] = BigIntModG{}

type BigIntModG struct {
	value *big.Int
	mod   *big.Int
}

func mod(x *big.Int, y *big.Int) *big.Int {
	return (&big.Int{}).Mod(x, y)
}

func add(x *big.Int, y *big.Int) *big.Int {
	return (&big.Int{}).Add(x, y)
}

func sub(x *big.Int, y *big.Int) *big.Int {
	return (&big.Int{}).Sub(x, y)
}

func addMod(x *big.Int, y *big.Int, m *big.Int) *big.Int {
	return mod(add(mod(x, m), mod(y, m)), m)
}

func subMod(x *big.Int, y *big.Int, m *big.Int) *big.Int {
	return mod(sub(mod(x, m), mod(y, m)), m)
}

func FromInt64(v int64, m int64) group.Elem[BigIntModG] {
	return BigIntModG{
		value: big.NewInt(v),
		mod:   big.NewInt(m),
	}
}

func FromBigInt(v *big.Int, m *big.Int) group.Elem[BigIntModG] {
	return BigIntModG{
		value: v,
		mod:   m,
	}
}

func (x BigIntModG) CombineWith(y group.Elem[BigIntModG]) group.Elem[BigIntModG] {
	return BigIntModG{
		value: addMod(x.value, y.AsPure().value, x.mod),
		mod:   x.mod,
	}
}

func (x BigIntModG) Invert() group.Elem[BigIntModG] {
	return BigIntModG{
		value: subMod(x.mod, x.value, x.mod),
		mod:   x.mod,
	}
}

func (x BigIntModG) EqualsTo(y group.Elem[BigIntModG]) bool {
	return x.value.Cmp(y.AsPure().value) == 0
}

func (x BigIntModG) AsPure() BigIntModG {
	return x
}

func (x BigIntModG) AsGroup() group.Elem[BigIntModG] {
	return x
}
