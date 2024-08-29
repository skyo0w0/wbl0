package cache

import (
	"sync"
)

type Cache struct {
	store sync.Map
}

func New() *Cache {
	return &Cache{}
}

func (c *Cache) Set(key string, value interface{}) {
	c.store.Store(key, value)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.store.Load(key)
}
