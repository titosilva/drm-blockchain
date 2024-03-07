package di

import (
	"fmt"
	"reflect"
)

type DIContext struct {
	Instances map[reflect.Type]any
	Factories map[reflect.Type]any
}

func ensureInitialized(provider *DIContext) {
	if provider.Instances == nil {
		provider.Instances = make(map[reflect.Type]any)
	}

	if provider.Factories == nil {
		provider.Factories = make(map[reflect.Type]any)
	}
}

func AddSingleton[T any](provider *DIContext, instance *T) {
	ensureInitialized(provider)
	tType := reflect.TypeOf([0]T{}).Elem()
	provider.Instances[tType] = instance
}

func AddFactory[T any](provider *DIContext, factory func() *T) {
	ensureInitialized(provider)
	tType := reflect.TypeOf([0]T{}).Elem()
	provider.Factories[tType] = factory
}

func GetService[T any](provider *DIContext) *T {
	ensureInitialized(provider)
	tType := reflect.TypeOf([0]T{}).Elem()
	serviceGeneric, ok := provider.Instances[tType]

	if !ok {
		serviceFactory, ok := provider.Factories[tType]

		if !ok {
			str := fmt.Sprintf("Service %s singleton or factory not found", tType.Name())
			panic(str)
		}

		serviceGeneric = serviceFactory.(func() *T)()
	}

	service, ok := serviceGeneric.(*T)
	if !ok {
		str := fmt.Sprintf("Failed to convert service %s", tType.Name())
		panic(str)
	}
	return service
}
