package queryable

import (
	"drm-blockchain/src/collections/iterator"
	"drm-blockchain/src/utils/maybe"
)

type Predicate[T any] func(T) bool

type Queryable[T any] interface {
	iterator.Iterable[T]
	All(pred Predicate[T]) bool
	Any(pred Predicate[T]) bool
	Where(pred Predicate[T]) Queryable[T]
	First() maybe.Maybe[T]

	Skip(offset int) Queryable[T]
	Take(count int) Queryable[T]
}
