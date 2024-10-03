package waitmap

import (
	"context"
	"sync"
)

// WaitMap is a key-value store that enables not only setting and getting
// values from a map, but also waiting until value for a key becomes available.
// Important: Since this implementation is tied pretty closely to how it will be used,
// (as an internal package), it has the peculiar behavior of only the first `Set` setting
// the value. Subsequent `Sets()` are recorded, but don't change the returned value.
type WaitMap[K comparable, V any] struct {
	sync.Mutex

	m    map[K][]V
	subs map[K][]chan V
}

// New creates a new WaitMap.
func New[K comparable, V any]() *WaitMap[K, V] {

	wm := WaitMap[K, V]{
		m:    make(map[K][]V),
		subs: make(map[K][]chan V),
	}

	return &wm
}

// Set sets the value for a key. If the value already exists, we append it to a list.
func (w *WaitMap[K, V]) Set(key K, value V) {
	w.Lock()
	defer w.Unlock()

	_, ok := w.m[key]
	if !ok {
		w.m[key] = make([]V, 0)
	}

	w.m[key] = append(w.m[key], value)

	// Send the new value to any waiting subscribers of the key.
	for _, sub := range w.subs[key] {
		sub <- value
	}
	delete(w.subs, key)
}

// Wait will wait until the value for a key becomes available.
func (w *WaitMap[K, V]) Wait(key K) V {
	w.Lock()
	// Unlock cannot be deferred so we can ublock Set() while waiting.

	values, ok := w.m[key]
	if ok {
		w.Unlock()
		return values[0]
	}

	// If there's no value yet, subscribe to any new values for this key.
	ch := make(chan V)
	w.subs[key] = append(w.subs[key], ch)
	w.Unlock()

	return <-ch
}

// WaitFor will wait for the value for a key to become available, but no longer than the specified duration.
func (w *WaitMap[K, V]) WaitFor(ctx context.Context, key K) (V, bool) {
	w.Lock()
	// Unlock cannot be deferred so we can ublock Set() while waiting.

	values, ok := w.m[key]
	if ok {
		w.Unlock()
		return values[0], true
	}

	// If there's no value yet, subscribe to any new values for this key.
	// Use a bufferred channel since we might bail before collecting our value.
	ch := make(chan V, 1)
	w.subs[key] = append(w.subs[key], ch)
	w.Unlock()

	select {
	case <-ctx.Done():
		zero := *new(V)
		return zero, false
	case value := <-ch:
		return value, true
	}
}

// Get will return the current value for the key, if any.
func (w *WaitMap[K, V]) Get(key K) (V, bool) {
	w.Lock()
	defer w.Unlock()

	values, ok := w.m[key]
	if !ok {
		zero := *new(V)
		return zero, ok
	}

	// As noted in the comment at the beginning of this file,
	// this is special behavior because of the way this map will be used.
	// Get will always return the first value.
	value := values[0]
	return value, true
}
