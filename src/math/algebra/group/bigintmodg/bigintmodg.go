package bigintmodg

import (
	"drm-blockchain/src/math/algebra/group"
	"math/big"
)

// static interface check
var _ group.Elem[BigIntModG] = BigIntModG{}

type BigIntModG struct {
	Value *big.Int
	Mod   *big.Int
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
	bmod := big.NewInt(m)
	return BigIntModG{
		Value: mod(big.NewInt(v), bmod),
		Mod:   bmod,
	}
}

func FromBigInt(v *big.Int, m *big.Int) group.Elem[BigIntModG] {
	return BigIntModG{
		Value: v,
		Mod:   m,
	}
}

func (x BigIntModG) CombineWith(y group.Elem[BigIntModG]) group.Elem[BigIntModG] {
	return BigIntModG{
		Value: addMod(x.Value, y.AsPure().Value, x.Mod),
		Mod:   x.Mod,
	}
}

func (x BigIntModG) Invert() group.Elem[BigIntModG] {
	return BigIntModG{
		Value: subMod(x.Mod, x.Value, x.Mod),
		Mod:   x.Mod,
	}
}

func (x BigIntModG) EqualsTo(y group.Elem[BigIntModG]) bool {
	return x.Value.Cmp(y.AsPure().Value) == 0
}

func (x BigIntModG) AsPure() BigIntModG {
	return x
}

func (x BigIntModG) AsGroup() group.Elem[BigIntModG] {
	return x
}
