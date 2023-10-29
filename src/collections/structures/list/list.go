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

func NewWithSize[T any](size int, fill func(int) T) List[T] {
	l := List[T]{
		values: make([]T, 0),
	}

	for i := 0; i < size; i++ {
		l.values = append(l.values, fill(i))
	}

	return l
}

func New[T any]() List[T] {
	return List[T]{
		values: make([]T, 0),
	}
}

func NewFrom[T any](elems []T) List[T] {
	return List[T]{
		values: elems,
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

func (l List[T]) Skip(offset int) queryable.Queryable[T] {
	if offset >= l.Count() {
		return NewFrom(make([]T, 0))
	}

	return NewFrom(l.values[offset:l.Count()])
}

func (l List[T]) Take(count int) queryable.Queryable[T] {
	if count >= l.Count() {
		return l
	}

	return NewFrom(l.values[0:count])
}
