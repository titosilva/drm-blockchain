package di

import (
	"fmt"
	"reflect"
)

type DIContext struct {
	SingletonInstances map[string]any
	SingletonFactories map[string]any
	Factories          map[string]any

	InterfaceSingletonInstances map[string]any
	InterfaceSingletonFactories map[string]any
	InterfaceFactories          map[string]any
}

func NewContext() *DIContext {
	ctx := new(DIContext)
	ctx.SingletonInstances = make(map[string]any)
	ctx.SingletonFactories = make(map[string]any)
	ctx.Factories = make(map[string]any)

	ctx.InterfaceSingletonInstances = make(map[string]any)
	ctx.InterfaceSingletonFactories = make(map[string]any)
	ctx.InterfaceFactories = make(map[string]any)
	return ctx
}

func AddSingleton[T any](provider *DIContext, factory func(*DIContext) *T) {
	tName := getTypeName[T]()
	provider.SingletonFactories[tName] = factory
}

func AddInterfaceSingleton[T any](provider *DIContext, factory func(*DIContext) T) {
	tName := getTypeName[T]()
	provider.InterfaceSingletonFactories[tName] = factory
}

func AddFactory[T any](provider *DIContext, factory func(*DIContext) *T) {
	tName := getTypeName[T]()
	provider.Factories[tName] = factory
}

func AddInterfaceFactory[T any](provider *DIContext, factory func(*DIContext) T) {
	tName := getTypeName[T]()
	provider.InterfaceFactories[tName] = factory
}

func GetService[T any](provider *DIContext) *T {
	tName := getTypeName[T]()

	tGeneric, found := provider.SingletonInstances[tName]
	if found {
		t, ok := tGeneric.(*T)

		if ok {
			return t
		}
	}

	serviceFactory, found := provider.SingletonFactories[tName]
	if found {
		tFactory, ok := serviceFactory.(func(*DIContext) *T)

		if ok {
			t := tFactory(provider)
			provider.SingletonInstances[tName] = t
			return t
		}
	}

	serviceFactory, found = provider.Factories[tName]
	if found {
		tFactory, ok := serviceFactory.(func(*DIContext) *T)

		if ok {
			return tFactory(provider)
		}
	}

	var panicMsg string
	if !found {
		panicMsg = fmt.Sprintf("Service %s singleton or factory not found", tName)
	} else {
		panicMsg = fmt.Sprintf("Failed to convert service %s", tName)
	}
	panic(panicMsg)
}

func GetInterfaceService[T any](provider *DIContext) T {
	tName := getTypeName[T]()

	tGeneric, found := provider.InterfaceSingletonInstances[tName]
	if found {
		t, ok := tGeneric.(T)

		if ok {
			return t
		}
	}

	serviceFactory, found := provider.InterfaceSingletonFactories[tName]
	if found {
		tFactory, ok := serviceFactory.(func(*DIContext) T)

		if ok {
			t := tFactory(provider)
			provider.InterfaceSingletonInstances[tName] = t
			return t
		}
	}

	serviceFactory, found = provider.InterfaceFactories[tName]
	if found {
		tFactory, ok := serviceFactory.(func(*DIContext) T)

		if ok {
			return tFactory(provider)
		}
	}

	var panicMsg string
	if !found {
		panicMsg = fmt.Sprintf("Interface %s singleton or factory not found", tName)
	} else {
		panicMsg = fmt.Sprintf("Interface to convert service %s", tName)
	}
	panic(panicMsg)
}

func getTypeName[T any]() string {
	tType := reflect.TypeOf([0]T{}).Elem()
	return tType.PkgPath() + "/" + tType.Name()
}
