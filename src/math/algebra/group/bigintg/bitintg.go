package bigintg

import (
	"drm-blockchain/src/math/algebra/group"
	"math/big"
)

// static interface check
var _ group.Elem[BigIntG] = BigIntG{}

type BigIntG struct {
	value *big.Int
}

func From(v int64) group.Elem[BigIntG] {
	return BigIntG{value: big.NewInt(v)}
}

func (x BigIntG) CombineWith(y group.Elem[BigIntG]) group.Elem[BigIntG] {
	return BigIntG{
		value: (&big.Int{}).Add(x.value, y.AsPure().value),
	}
}

func (x BigIntG) Zero() group.Elem[BigIntG] {
	return From(0)
}

func (x BigIntG) EqualsTo(y group.Elem[BigIntG]) bool {
	return x.value.Cmp(y.AsPure().value) == 0
}

func (x BigIntG) Invert() group.Elem[BigIntG] {
	return BigIntG{value: (&big.Int{}).Mul(x.value, big.NewInt(-1))}
}

func (x BigIntG) AsPure() BigIntG {
	return x
}

func (x BigIntG) AsGroup() group.Elem[BigIntG] {
	return x
}
