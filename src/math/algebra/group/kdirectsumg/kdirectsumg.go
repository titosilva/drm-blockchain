package kdirectsumg

import (
	"drm-blockchain/src/collections/structures/list"
	"drm-blockchain/src/math/algebra/group"
	"drm-blockchain/src/math/algebra/group/bigintmodg"
)

// static interface check
var _ group.Elem[KDirectSumG[bigintmodg.BigIntModG]] = KDirectSumG[bigintmodg.BigIntModG]{}

type KDirectSumG[X any] struct {
	k       int
	Entries list.List[group.Elem[X]]
}

func New[X any](k int) KDirectSumG[X] {
	return KDirectSumG[X]{
		k:       k,
		Entries: list.NewWithSize[group.Elem[X]](k),
	}
}

func (x KDirectSumG[G]) CombineWith(y group.Elem[KDirectSumG[G]]) group.Elem[KDirectSumG[G]] {
	if x.k != y.AsPure().k {
		panic("Can't combine KDirectSum with different k")
	}

	r := New[G](x.k)
	r_iter := r.Entries.GetIterator()
	x_iter := x.Entries.GetIterator()
	y_iter := y.AsPure().Entries.GetIterator()

	for x_iter.HasNext() && y_iter.HasNext() && r_iter.HasNext() {
		xi := *(x_iter.GetNext())
		yi := *(y_iter.GetNext())
		rip := r_iter.GetNext()

		*rip = xi.CombineWith(yi)
	}

	return r
}

func (x KDirectSumG[G]) Invert() group.Elem[KDirectSumG[G]] {
	r := New[G](x.k)
	x_iter := x.Entries.GetIterator()

	for x_iter.HasNext() {
		xi := *(x_iter.GetNext())
		r.Entries.Add(xi.Invert())
	}

	return r
}

func (x KDirectSumG[G]) EqualsTo(y group.Elem[KDirectSumG[G]]) bool {
	if x.k != y.AsPure().k {
		panic("Can't compare KDirectSum with different k")
	}

	x_iter := x.Entries.GetIterator()
	y_iter := y.AsPure().Entries.GetIterator()

	for x_iter.HasNext() && y_iter.HasNext() {
		xi := *(x_iter.GetNext())
		yi := *(y_iter.GetNext())

		if !xi.EqualsTo(yi) {
			return false
		}
	}

	return true
}

func (x KDirectSumG[G]) AsPure() KDirectSumG[G] {
	return x
}

func (x KDirectSumG[G]) AsGroup() group.Elem[KDirectSumG[G]] {
	return x
}
