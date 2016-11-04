package main

import (
	"log"
	"net/http"

	"github.com/satori/go.uuid"
	"github.com/yanpozka/checkers/store"
)

func createGame(w http.ResponseWriter, r *http.Request) {
	ms, isStoreType := r.Context().Value(storeCtxKey).(store.Store)
	if !isStoreType {
		log.Fatal("This's very fatal, Context doesn't have storeCtxKey :(")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data := []byte("abc") // SUPER TODO: serialize a Game obj
	key := "game-" + uuid.NewV4().String()

	if err := ms.Set(key, data); err != nil {
		log.Println("Error trying to create (Set) new game", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("true"))
}
