package store

import (
	"errors"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

// ErrNotFoundItem custom error to know if a key doesn't exist
var ErrNotFoundItem = errors.New("key not found")

// Store basic interface for key/value store
type Store interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
}

type mem struct {
	c *memcache.Client
}

// NewMemcacheStore creates new memcache Store
func NewMemcacheStore(s ...string) Store {
	mc := memcache.New(s...)

	return &mem{c: mc}
}

func (m *mem) Get(key string) ([]byte, error) {

	item, err := m.c.Get(key)

	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, ErrNotFoundItem
		}
		return nil, fmt.Errorf("memcahe error: %v\n", err)
	}

	return item.Value, nil
}

func (m *mem) Set(key string, val []byte) error {

	err := m.c.Set(&memcache.Item{Key: key, Value: val})

	if err != nil {
		return fmt.Errorf("memcahe error: %v\n", err)
	}

	return nil
}
