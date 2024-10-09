package waitmap

import (
	"context"
	"math"
	"sync"

	"github.com/hashicorp/golang-lru/simplelru"
)

// WaitMap is a key-value store that enables not only setting and getting
// values from a map, but also waiting until value for a key becomes available.
type WaitMap[K comparable, V any] struct {
	sync.Mutex

	cache *simplelru.LRU
	subs  map[K][]chan V
}

// New creates a new WaitMap.
func New[K comparable, V any](size int) *WaitMap[K, V] {

	if size <= 0 {
		size = math.MaxInt
	}

	// Only possible cause of an error is providing an invalid size value
	cache, _ := simplelru.NewLRU(size, nil)

	wm := WaitMap[K, V]{
		cache: cache,
		subs:  make(map[K][]chan V),
	}

	return &wm
}

// Set sets the value for a key. If the value already exists, we append it to a list.
func (w *WaitMap[K, V]) Set(key K, value V) {
	w.Lock()
	defer w.Unlock()

	w.cache.Add(key, value)

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

	value, ok := w.cache.Get(key)
	if ok {
		w.Unlock()
		return value.(V)
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

	value, ok := w.cache.Get(key)
	if ok {
		w.Unlock()
		return value.(V), true
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

	value, ok := w.cache.Get(key)
	if !ok {
		zero := *new(V)
		return zero, ok
	}

	return value.(V), true
}
