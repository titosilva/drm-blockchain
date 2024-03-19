package di

import (
	"fmt"
	"reflect"
)

type DIContext struct {
	Instances          map[reflect.Type]any
	Factories          map[reflect.Type]any
	InterfaceInstances map[reflect.Type]any
	InterfaceFactories map[reflect.Type]any
}

func NewContext() *DIContext {
	ctx := new(DIContext)
	ctx.Instances = make(map[reflect.Type]any)
	ctx.Factories = make(map[reflect.Type]any)
	return ctx
}

func ensureInitialized(provider *DIContext) {
	if provider.Instances == nil {
		provider.Instances = make(map[reflect.Type]any)
	}

	if provider.InterfaceInstances == nil {
		provider.InterfaceInstances = make(map[reflect.Type]any)
	}

	if provider.Factories == nil {
		provider.Factories = make(map[reflect.Type]any)
	}

	if provider.InterfaceFactories == nil {
		provider.InterfaceFactories = make(map[reflect.Type]any)
	}
}

func AddSingleton[T any](provider *DIContext, instance *T) {
	ensureInitialized(provider)
	tType := reflect.TypeOf([0]T{}).Elem()
	provider.Instances[tType] = instance
}

func AddInterfaceSingleton[T any](provider *DIContext, instance T) {
	ensureInitialized(provider)
	tType := reflect.TypeOf([0]T{}).Elem()
	provider.InterfaceInstances[tType] = instance
}

func AddFactory[T any](provider *DIContext, factory func(*DIContext) *T) {
	ensureInitialized(provider)
	tType := reflect.TypeOf([0]T{}).Elem()
	provider.Factories[tType] = factory
}

func AddInterfaceFactory[T any](provider *DIContext, factory func(*DIContext) T) {
	ensureInitialized(provider)
	tType := reflect.TypeOf([0]T{}).Elem()
	provider.InterfaceFactories[tType] = factory
}

func GetService[T any](provider *DIContext) *T {
	ensureInitialized(provider)
	tType := reflect.TypeOf([0]T{}).Elem()

	serviceGeneric, found := provider.Instances[tType]
	if found {
		if service, ok := serviceGeneric.(*T); ok {
			return service
		}
	}

	serviceFactory, found := provider.Factories[tType]
	if found {
		serviceGeneric = serviceFactory.(func(*DIContext) *T)(provider)
		if service, ok := serviceGeneric.(*T); ok {
			return service
		}
	}

	var panicMsg string
	if !found {
		panicMsg = fmt.Sprintf("Service %s singleton or factory not found", tType.Name())
	} else {
		panicMsg = fmt.Sprintf("Failed to convert service %s", tType.Name())
	}
	panic(panicMsg)
}

func GetInterfaceService[T any](provider *DIContext) T {
	tType := reflect.TypeOf([0]T{}).Elem()

	serviceGeneric, found := provider.InterfaceInstances[tType]
	if found {
		if service, ok := serviceGeneric.(T); ok {
			return service
		}
	}

	serviceFactory, found := provider.InterfaceFactories[tType]
	if found {
		serviceGeneric = serviceFactory.(func(*DIContext) T)(provider)
		if service, ok := serviceGeneric.(T); ok {
			return service
		}
	}

	var panicMsg string
	if !found {
		panicMsg = fmt.Sprintf("Service %s singleton or factory not found", tType.Name())
	} else {
		panicMsg = fmt.Sprintf("Failed to convert service %s", tType.Name())
	}
	panic(panicMsg)
}
