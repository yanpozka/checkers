package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/yanpozka/checkers/store"
)

const storeCtxKey = "store"

func loggerMW(inner http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf("%s: %q %s Time consumed: %s", r.Method, r.URL, r.RemoteAddr, time.Since(start))
	})
}

func commonMW(inner http.Handler) http.Handler {

	// TODO: it can be a list of hosts separated by comma
	ms := store.NewMemcacheStore(getOrDefault("MEMCACHE_PORT", defaultMemcacheHost))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Server", "GoYan")

		// default application/json for api responses for now
		header.Set("Content-Type", "application/json")

		ctx := context.WithValue(r.Context(), storeCtxKey, ms)

		inner.ServeHTTP(w, r.WithContext(ctx))
	})
}
