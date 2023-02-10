package cache

import (
	"sync"

	"github.com/blocklessnetworking/b7s/models/execute"
)

// Cache is a simple cache storing execution responses in memory.
type Cache struct {
	sync.Mutex
	m map[string]*execute.Response
}

// New creates a new cache.
func New() *Cache {
	c := Cache{
		m: make(map[string]*execute.Response),
	}
	return &c
}

// Get retrieves an execution response from the cache, given its executionID.
func (c *Cache) Get(id string) (*execute.Response, bool) {
	c.Lock()
	defer c.Unlock()

	res, ok := c.m[id]
	return res, ok
}

// Set caches the given execution response.
func (c *Cache) Set(id string, res *execute.Response) {
	c.Lock()
	defer c.Unlock()

	c.m[id] = res
}
