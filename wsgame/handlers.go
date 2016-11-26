package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/yanpozka/checkers/store"
)

func gameWS(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	gameID := parts[len(parts)-1]
	if gameID == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	playerID := r.URL.Query().Get("player")
	if playerID == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if bg, err := ms.Get(playerID); err == store.ErrItemNotFound || string(bg) != gameID {
		log.Printf("PlayerID: %q doesn't have an associated gameID: %q\n", playerID, gameID)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil) // if Upgrade fails, it'll write an error message
	if err != nil {
		log.Println("Error trying to Upgrade: ", err)
		return
	}

	c := client{
		conn:     conn,
		send:     make(chan []byte, 1),
		done:     make(chan struct{}, 1),
		gameID:   gameID,
		playerID: playerID,
	}

	go c.write()

	if err := c.startListenFromNats(); err != nil {
		log.Println("Error trying to subscribe:", err)
		c.done <- struct{}{}
		return
	}

	c.read()
}
