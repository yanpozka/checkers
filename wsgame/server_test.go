package main

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats"
)

const delayms = 60 // 60 FPS ? :P

var tserver *httptest.Server

//
func TestServerOK(t *testing.T) {
	tserver = httptest.NewServer(createServer())

	gameID := "game-123"
	u := url.URL{Scheme: "ws", Host: getHost(tserver.URL), Path: "/game/" + gameID, RawQuery: "player=player-123"} // ?
	t.Logf("Connecting to %q", u.String())

	c, res, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatal("Error on dial:", err, "Status:", res.Status)
	}
	defer c.Close()

	sendAndRecv(t, c)

	nc, err := nats.Connect(getOrDefault("NATS_URL", nats.DefaultURL))
	if err != nil {
		t.Fatal(err)
	}
	defer nc.Close()

	subj, msg := "subject-"+gameID, []byte("only rest check the logs :( so far")
	nc.Publish(subj, msg)

	time.Sleep(delayms*time.Millisecond + 1)
}

//
func sendAndRecv(t *testing.T, ws *websocket.Conn) {
	if err := ws.SetWriteDeadline(time.Now().Add(delayms * time.Millisecond)); err != nil {
		t.Fatalf("SetWriteDeadline: %v", err)
	}
	if err := ws.SetReadDeadline(time.Now().Add(delayms * time.Millisecond)); err != nil {
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
