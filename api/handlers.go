package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/satori/go.uuid"
	"github.com/yanpozka/checkers/game"
	"github.com/yanpozka/checkers/store"
)

func createGame(w http.ResponseWriter, r *http.Request) {
	ms, isStoreType := r.Context().Value(storeCtxKey).(store.Store)
	if !isStoreType {
		panic("This's very fatal, Context doesn't have storeCtxKey :(")
	}

	playerID := "player-" + uuid.NewV4().String()
	data := game.InitGame([]byte(playerID))

	gameID := "game-" + uuid.NewV4().String()

	if writeErr(ms.Set(gameID, data), w, "Error creating (Set) new game") {
		return
	}
	if writeErr(ms.Set(playerID, []byte(gameID)), w, "Error setting new game with player (author)") {
		return
	}

	mr := map[string]string{
		"playerID": playerID,
		"gameURL":  "/api/invitation/" + gameID,
	}

	buf, err := json.Marshal(mr)

	if !writeErr(err, w, "Encoding response") {
		w.WriteHeader(http.StatusCreated)
		w.Write(buf)
	}
}

func invitation(w http.ResponseWriter, r *http.Request) {
	_, isStoreType := r.Context().Value(storeCtxKey).(store.Store)
	if !isStoreType {
		panic("This's very fatal, Context doesn't have storeCtxKey :(")
	}
}

func writeErr(err error, w http.ResponseWriter, message string) bool {
	if err != nil {
		log.Println(message, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return true
	}
	return false
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("true"))
}
