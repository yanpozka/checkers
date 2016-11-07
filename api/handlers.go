package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/satori/go.uuid"
	"github.com/yanpozka/checkers/game"
	"github.com/yanpozka/checkers/store"
)

const (
	Empty byte = iota
	PlayerUp
	PlayerDown
)

func createGame(w http.ResponseWriter, r *http.Request) {
	ms, isStoreType := r.Context().Value(storeCtxKey).(store.Store)
	if !isStoreType {
		log.Fatal("This's very fatal, Context doesn't have storeCtxKey :(")
	}

	playerID := "player-" + uuid.NewV4().String()
	data := game.InitGame([]byte(playerID))

	gameID := "game-" + uuid.NewV4().String()

	if checkErr(ms.Set(gameID, data), w, "Error creating (Set) new game") {
		return
	}
	if checkErr(ms.Set(playerID, []byte(gameID)), w, "Error setting new game with player (author)") {
		return
	}

	w.WriteHeader(http.StatusCreated)

	mr := map[string]string{
		"playerID": playerID,
		"gameURL":  "/api/invitation/" + gameID,
	}

	checkErr(json.NewEncoder(w).Encode(&mr), w, "")
}

func checkErr(err error, w http.ResponseWriter, message string) bool {
	if err != nil {
		log.Println(message, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	return err != nil
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("true"))
}
