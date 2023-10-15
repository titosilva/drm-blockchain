package group

// Elem defines an interface for an element of an abstract group
type Elem[G any] interface {
	Zero() Elem[G]
	CombineWith(a Elem[G]) Elem[G]
	Invert() Elem[G]
	AsPure() G
	AsGroup() Elem[G]
}
