package kdirectsumg

import (
	"drm-blockchain/src/collections/structures/list"
	"drm-blockchain/src/math/algebra/group"
	"drm-blockchain/src/math/algebra/group/bigintmodg"
)

// static interface check
var _ group.Elem[KDirectSumG[bigintmodg.BigIntModG]] = KDirectSumG[bigintmodg.BigIntModG]{}

type KDirectSumG[X any] struct {
	k      int
	values list.List[group.Elem[X]]
}

func New[X any](k int) KDirectSumG[X] {
	return KDirectSumG[X]{
		k:      k,
		values: list.NewWithSize[group.Elem[X]](k),
	}
}

func (x KDirectSumG[G]) CombineWith(y group.Elem[KDirectSumG[G]]) group.Elem[KDirectSumG[G]] {
	if x.k != y.AsPure().k {
		panic("Can't combine KDirectSum with different k")
	}

	r := New[G](x.k)
	x_iter := x.values.GetIterator()
	y_iter := y.AsPure().values.GetIterator()

	for x_iter.HasNext() && y_iter.HasNext() {

	}

	return r
}

func (x KDirectSumG[G]) Invert() group.Elem[KDirectSumG[G]] {
	return KDirectSumG[G]{}
}

func (x KDirectSumG[G]) EqualsTo(y group.Elem[KDirectSumG[G]]) bool {
	return true
}

func (x KDirectSumG[G]) AsPure() KDirectSumG[G] {
	return x
}

func (x KDirectSumG[G]) AsGroup() group.Elem[KDirectSumG[G]] {
	return x
}
