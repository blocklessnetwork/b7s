package cache

import (
	"sync"

	"github.com/blocklessnetworking/b7s/models/execute"
)

// Cache is a simple cache storing execution responses in memory.
type Cache struct {
	sync.Mutex
	m map[string]execute.Result
}

// New creates a new cache.
func New() *Cache {
	c := Cache{
		m: make(map[string]execute.Result),
	}
	return &c
}

// Get retrieves an execution response from the cache, given its requestID.
func (c *Cache) Get(id string) (execute.Result, bool) {
	c.Lock()
	defer c.Unlock()

	res, ok := c.m[id]
	return res, ok
}

// Set caches the given execution response.
func (c *Cache) Set(id string, res execute.Result) {
	c.Lock()
	defer c.Unlock()

	c.m[id] = res
}
