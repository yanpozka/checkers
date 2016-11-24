package store

import (
	"os"
	"strings"
	"testing"
)

func TestSetandGet(t *testing.T) {

	hosts := strings.Split(os.Getenv("MEMCACHE_PORT"), "")
	if len(hosts) == 0 {
		hosts = append(hosts, "127.0.0.1:11211")
	}

	store := NewMemcacheStore(hosts...)

	key := "mea_culpa"
	val := []byte("abcdefg")

	err := store.Set(key, val)
	if err != nil {
		t.Fatal(err)
	}

	v, err := store.Get(key)
	if err != nil {
		t.Fatal(err)
	}

	if string(v) != string(val) {
		t.Fatalf("Expected: %q, got: %q\n", string(val), string(v))
	}
}
