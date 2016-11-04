package store

import (
	"os"
	"testing"
)

func TestGetSet(t *testing.T) {
	store := NewMemcacheStore(os.Getenv("MEMCACHE_PORT"))

	key := "mea_culpa"
	val := []byte("abcd")

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
