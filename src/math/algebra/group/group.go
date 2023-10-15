package group

// Elem defines an interface for an element of an abstract group
type Elem[G any] interface {
	CombineWith(a Elem[G]) Elem[G]
	Invert() Elem[G]
	EqualsTo(a Elem[G]) bool
	AsPure() G
	AsGroup() Elem[G]
}
