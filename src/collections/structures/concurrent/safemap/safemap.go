package safemap

import "sync"

type SafeMap[TKey comparable, TVal any] struct {
	underlyingMap  map[TKey]TVal
	readWriteMutex *sync.RWMutex
}

func New[TKey comparable, TVal any]() *SafeMap[TKey, TVal] {
	m := new(SafeMap[TKey, TVal])

	m.underlyingMap = make(map[TKey]TVal)
	m.readWriteMutex = new(sync.RWMutex)

	return m
}

func (m *SafeMap[TKey, TVal]) Get(key TKey) (TVal, bool) {
	m.readWriteMutex.Lock()
	val, found := m.underlyingMap[key]
	m.readWriteMutex.Unlock()

	return val, found
}

func (m *SafeMap[TKey, TVal]) Set(key TKey, val TVal) {
	m.readWriteMutex.Lock()
	m.underlyingMap[key] = val
	m.readWriteMutex.Unlock()
}

func (m *SafeMap[TKey, TVal]) Delete(key TKey) {
	delete(m.underlyingMap, key)
}
