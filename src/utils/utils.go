package utils

import "reflect"

func TypeToString[T any]() string {
	return reflect.TypeOf([0]T{}).Elem().Name()
}

func TypeOf(t any) string {
	return reflect.TypeOf(t).Name()
}
