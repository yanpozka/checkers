package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"github.com/yanpozka/checkers/game"
	"github.com/yanpozka/checkers/store"
)

func createGame(w http.ResponseWriter, r *http.Request) {
	st, isStoreType := r.Context().Value(storeCtxKey).(store.Store)
	if !isStoreType {
		panic("This's very fatal, Context doesn't have storeCtxKey :(")
	}

	playerID := "player-" + uuid.NewV4().String()
	data := game.InitGame([]byte(playerID))

	gameID := "game-" + uuid.NewV4().String()

	if writeErr(st.Set(gameID, data), w, "Error creating (Set) new game") {
		return
	}
	if writeErr(st.Set(playerID, []byte(gameID)), w, "Error setting new game with player (author)") {
		return
	}

	mr := map[string]string{
		"playerID":      playerID,
		"invitationURL": "/api/invitation/" + gameID,
	}

	buf, err := json.Marshal(mr)

	if !writeErr(err, w, "Encoding response") {
		w.WriteHeader(http.StatusCreated)
		w.Write(buf)
	}
}

func invitation(w http.ResponseWriter, r *http.Request) {
	st, isStoreType := r.Context().Value(storeCtxKey).(store.Store)
	if !isStoreType {
		panic("Context doesn't have storeCtxKey or isn't store.Store type")
	}

	params, isParams := r.Context().Value(paramsCtxKey).(httprouter.Params)
	if !isParams {
		panic("Context doesn't have storeCtxKey or isn't Params type")
	}
	gameID := params.ByName("gameID")

	if _, err := st.Get(gameID); err == store.ErrItemNotFound {
		writeErr(err, w, "Game ID not found: "+gameID)
		return
	}

	playerID := "player-" + uuid.NewV4().String()
	if writeErr(ms.Set(playerID, []byte(gameID)), w, "Error setting second player on invitation") {
		return
	}

	mr := map[string]string{
		"playerID": playerID,
		"gameURL":  "/ws/game/" + gameID,
	}

	buf, err := json.Marshal(mr)

	if !writeErr(err, w, "Encoding response") {
		w.WriteHeader(http.StatusOK)
		w.Write(buf)
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
