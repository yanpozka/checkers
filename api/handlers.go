package main

import "net/http"

func createGame(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("true"))
}
