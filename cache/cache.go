package cache

import "sync"

type Cache interface {
	Set(k, v any)
	Get(k any) (any, bool)
}

type InMemoryCache struct {
	sync.Map
}

var AppCache = NewInMemoryCache()

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{}
}

func (c *InMemoryCache) Set(k, v any) {
	c.Store(k, v)
}

func (c *InMemoryCache) Get(k any) (any, bool) {
	return c.Load(k)
}
