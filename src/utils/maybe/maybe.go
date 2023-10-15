package maybe

import "errors"

type Maybe[T any] interface {
	Get() (T, error)
	IsNothing() bool
}

type just[T any] struct {
	value T
}

func Just[T any](val T) Maybe[T] {
	return just[T]{value: val}
}

func (j just[T]) Get() (T, error) {
	return j.value, nil
}

func (j just[T]) IsNothing() bool {
	return false
}

type nothing[T any] struct{}

func Nothing[T any]() Maybe[T] {
	return nothing[T]{}
}

func (n nothing[T]) Get() (T, error) {
	var a T
	return a, errors.New("cannot get from Nothing")
}

func (n nothing[T]) IsNothing() bool {
	return true
}
