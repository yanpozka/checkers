package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

func loggerMW(inner http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, req)

		log.Printf("%s: %q Time consumed: %s", r.Method, r.URL, r.RemoteAddr, time.Since(start))
	})
}

func commonMW(inner http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		header := w.Header()
		header.Set("Server", "GoYan")
		header.Set("X-Powered-By", "yandry-server 0.1")

		// default application/json for api responses for now
		header.Set("Content-Type", "application/json")

		inner.ServeHTTP(w, req)
	})
}
