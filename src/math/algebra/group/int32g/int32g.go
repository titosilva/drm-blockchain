package int32g

import "drm-blockchain/src/math/algebra/group"

// static interface check
var _ group.Elem[Int32G] = Int32G{}

type Int32G struct {
	value int32
}

func Zero() group.Elem[Int32G] {
	return Int32G{value: 0}
}

func (x Int32G) Zero() group.Elem[Int32G] {
	return Zero()
}

func From(v int32) group.Elem[Int32G] {
	return Int32G{value: v}
}

func (x Int32G) CombineWith(y group.Elem[Int32G]) group.Elem[Int32G] {
	return Int32G{
		value: x.value + y.AsPure().value,
	}
}

func (x Int32G) EqualsTo(y group.Elem[Int32G]) bool {
	return x == y
}

func (x Int32G) Invert() group.Elem[Int32G] {
	return Int32G{value: -1 * x.value}
}

func (x Int32G) AsPure() Int32G {
	return x
}

func (x Int32G) AsGroup() group.Elem[Int32G] {
	return x
}
