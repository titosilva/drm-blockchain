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

func New[X any](k int, fill func(i int) group.Elem[X]) KDirectSumG[X] {
	return KDirectSumG[X]{
		k:       k,
		Entries: list.NewWithSize[group.Elem[X]](k, fill),
	}
}

func (x KDirectSumG[G]) CombineWith(y group.Elem[KDirectSumG[G]]) group.Elem[KDirectSumG[G]] {
	if x.k != y.AsPure().k {
		panic("Can't combine KDirectSum with different k")
	}

	r := New[G](x.k, func(i int) group.Elem[G] { return x.Entries.ToArray()[0].Zero() })
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
	r := New[G](x.k, func(i int) group.Elem[G] { return x.Entries.ToArray()[0].Zero() })
	x_iter := x.Entries.GetIterator()
	r_iter := r.Entries.GetIterator()

	for x_iter.HasNext() && r_iter.HasNext() {
		xi := *(x_iter.GetNext())
		ri := r_iter.GetNext()

		*ri = xi.Invert()
	}

	return r
}

func (x KDirectSumG[G]) Zero() group.Elem[KDirectSumG[G]] {
	z := New[G](x.k, func(i int) group.Elem[G] { return x.Entries.ToArray()[0].Zero() })

	z_iter := z.Entries.GetIterator()

	for z_iter.HasNext() {
		*(z_iter.GetNext()) = x.AsPure().Entries.ToArray()[0].Zero()
	}

	return z
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
