package iterator

type Iterator[T any] interface {
	HasNext() bool
	GetNext() *T
}

type Iterable[T any] interface {
	GetIterator() Iterator[T]
	ToArray() []T
	Count() int
}
