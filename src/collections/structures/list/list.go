package list

import (
	"drm-blockchain/src/collections/iterator"
	"drm-blockchain/src/collections/queryable"
	"drm-blockchain/src/utils/maybe"
)

var _ iterator.Iterable[int] = List[int]{}
var _ queryable.Queryable[int] = List[int]{}
var _ iterator.Iterator[int] = &ListIterator[int]{}

type ListIterator[T any] struct {
	index int
	list  *List[T]
}

type List[T any] struct {
	values []T
}

func (l *List[T]) Add(elem T) {
	l.values = append(l.values, elem)
}

func NewWithSize[T any](size int) List[T] {
	return List[T]{
		values: make([]T, size),
	}
}

func New[T any]() List[T] {
	return List[T]{
		values: make([]T, 0),
	}
}

// ListIterator implements iterator.Iterator
func (li *ListIterator[T]) GetNext() *T {
	r := &li.list.values[li.index]
	li.index++
	return r
}

func (li ListIterator[T]) HasNext() bool {
	return len(li.list.values) > li.index
}

// List implements iterator.Iterable
func (l List[T]) Count() int {
	return len(l.values)
}

func (l List[T]) GetIterator() iterator.Iterator[T] {
	return &ListIterator[T]{
		index: 0,
		list:  &l,
	}
}

func (l List[T]) ToArray() []T {
	return l.values
}

// List implements iterator.Queryable
func (l List[T]) All(pred queryable.Predicate[T]) bool {
	for _, elem := range l.values {
		if !pred(elem) {
			return false
		}
	}

	return true
}

func (l List[T]) Any(pred queryable.Predicate[T]) bool {
	for _, elem := range l.values {
		if pred(elem) {
			return true
		}
	}

	return false
}

func (l List[T]) First() maybe.Maybe[T] {
	if len(l.values) == 0 {
		return maybe.Nothing[T]()
	}

	return maybe.Just(l.values[0])
}

func (l List[T]) Where(pred queryable.Predicate[T]) queryable.Queryable[T] {
	r := New[T]()

	for _, elem := range l.values {
		if pred(elem) {
			r.Add(elem)
		}
	}

	return r
}
