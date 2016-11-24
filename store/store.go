package store

import (
	"errors"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

//
var ErrNotFoundItem = errors.New("key not found")

//
type Store interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
}

type mem struct {
	c *memcache.Client
}

//
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
