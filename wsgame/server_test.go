package main

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var tserver *httptest.Server

func TestServerOK(t *testing.T) {
	tserver = httptest.NewServer(createServer())

	u := url.URL{Scheme: "ws", Host: getHost(tserver.URL), Path: "/game/game-123", RawQuery: "player=player-123"} // ?
	t.Logf("Connecting to %q", u.String())

	c, res, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatal("Error on dial:", err, "Status:", res.Status)
	}
	defer c.Close()

	sendAndRecv(t, c)
}

//
func sendAndRecv(t *testing.T, ws *websocket.Conn) {
	if err := ws.SetWriteDeadline(time.Now().Add(time.Millisecond * 500)); err != nil {
		t.Fatalf("SetWriteDeadline: %v", err)
	}
	if err := ws.SetReadDeadline(time.Now().Add(time.Millisecond * 500)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("Error Read Message: %v", err)
		}
		if !strings.Contains(string(p), "OK") {
			t.Fatalf("message=%s, want %s", p, "OK")
		}
		t.Log("Received=", string(p))
	}()

	if err := ws.WriteMessage(websocket.TextMessage, []byte("OK OK")); err != nil {
		t.Fatalf("WriteMessage: %v", err)
	}

	wg.Wait()
}

func getHost(url string) string {
	return strings.TrimPrefix(url, "http://")
}
