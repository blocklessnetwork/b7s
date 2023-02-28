package waitmap

import (
	"sync"
	"time"
)

// NOTE: Perhaps enable an option to say how long to wait for?

// WaitMap is a key-value store that enables not only setting and getting
// values from a map, but also waiting until value for a key becomes available.
// Important: Since this implementation is tied pretty closely to how it will be used,
// (as an internal package), it has the peculiar behavior of only the first `Set` setting
// the value. Subsequent `Sets()` are recorded, but don't change the returned value.
type WaitMap struct {
	sync.Mutex

	m    map[string][]any
	subs map[string][]chan any
}

// New creates a new WaitMap.
func New() *WaitMap {

	wm := WaitMap{
		m:    make(map[string][]any),
		subs: make(map[string][]chan any),
	}

	return &wm
}

// Set sets the value for a key. If the value already exists, we append it to a list.
func (w *WaitMap) Set(key string, value any) {
	w.Lock()
	defer w.Unlock()

	_, ok := w.m[key]
	if !ok {
		w.m[key] = make([]any, 0)
	}

	w.m[key] = append(w.m[key], value)

	// Send the new value to any waiting subscribers of the key.
	for _, sub := range w.subs[key] {
		sub <- value
	}
	delete(w.subs, key)
}

// Wait will wait until the value for a key becomes available.
func (w *WaitMap) Wait(key string) any {
	w.Lock()
	// Unlock cannot be deferred so we can ublock Set() while waiting.

	values, ok := w.m[key]
	if ok {
		w.Unlock()
		return values[0]
	}

	// If there's no value yet, subscribe to any new values for this key.
	ch := make(chan any)
	w.subs[key] = append(w.subs[key], ch)
	w.Unlock()

	return <-ch
}

// WaitFor will wait for the value for a key to become available, but no longer than the specified duration.
func (w *WaitMap) WaitFor(key string, d time.Duration) (any, bool) {
	w.Lock()
	// Unlock cannot be deferred so we can ublock Set() while waiting.

	values, ok := w.m[key]
	if ok {
		w.Unlock()
		return values[0], true
	}

	// If there's no value yet, subscribe to any new values for this key.
	// Use a bufferred channel since we might bail before collecting our value.
	ch := make(chan any, 1)
	w.subs[key] = append(w.subs[key], ch)
	w.Unlock()

	ticker := time.NewTicker(d)
	select {
	case <-ticker.C:
		return nil, false
	case value := <-ch:
		return value, true
	}
}

// Get will return the current value for the key, if any.
func (w *WaitMap) Get(key string) (any, bool) {
	w.Lock()
	defer w.Unlock()

	values, ok := w.m[key]
	if !ok {
		return values, ok
	}

	// As noted in the comment at the beginning of this file,
	// this is special behavior because of the way this map will be used.
	// Get will always return the first value.
	value := values[0]
	return value, true
}
