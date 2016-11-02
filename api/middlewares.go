package main

import (
	"log"
	"net/http"
	"time"
)

//
func loggerMW(inner http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ruri := req.URL.String()
		start := time.Now()

		inner.ServeHTTP(w, req) // !!

		log.Printf("%s: [ %s ] Time consumed: %s", req.Method, ruri, time.Since(start))
	})
}

//
func commonHeadersMW(inner http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		header := w.Header()
		header.Set("Server", "GoYan")
		header.Set("X-Powered-By", "yandry-server 0.1")

		// default application/json for api responses for now
		header.Set("Content-Type", "application/json")

		inner.ServeHTTP(w, req) // !!
	})
}

//

//
