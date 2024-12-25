package syncmap

import (
	"sync"
)

// NOTE: This package could be named "sync" to have a cleaner use of "sync.Map",
// but in order to not cause confusion with the std "sync.Map" we used this even though it stutters.

type Map[K comparable, V any] struct {
	*sync.RWMutex
	data map[K]V
}

func New[K comparable, V any]() *Map[K, V] {

	m := Map[K, V]{
		RWMutex: &sync.RWMutex{},
		data:    make(map[K]V),
	}

	return &m
}

func (m *Map[K, V]) Set(key K, value V) {
	m.Lock()
	defer m.Unlock()

	m.data[key] = value
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	m.RLock()
	defer m.RUnlock()

	val, ok := m.data[key]
	return val, ok
}

func (m *Map[K, V]) Delete(key K) {
	m.Lock()
	defer m.Unlock()

	delete(m.data, key)
}

func (m *Map[K, V]) Keys() []K {
	m.RLock()
	defer m.RUnlock()

	i := 0
	keys := make([]K, len(m.data))
	for key := range m.data {
		keys[i] = key
		i++
	}

	return keys
}

func (m *Map[K, V]) WithRLock(fn func(map[K]V)) {
	m.RLock()
	defer m.RUnlock()

	fn(m.data)
}

func (m *Map[K, V]) WithLock(fn func(map[K]V)) {
	m.Lock()
	defer m.Unlock()

	fn(m.data)
}
