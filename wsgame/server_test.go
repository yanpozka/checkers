package main

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var tserver *httptest.Server

func TestServerOK(t *testing.T) {
	tserver = httptest.NewServer(createServer())

	var dialer = websocket.Dialer{}

	url := makeWsProto(tserver.URL) + "/game"
	wscnn, resp, err := dialer.Dial(url, nil)

	if err != nil {
		t.Fatalf("Error trying to Dial() Resp='%v' Error='%v'", resp, err)
	}
	defer wscnn.Close()

	sendAndRecv(t, wscnn)
}

//
func sendAndRecv(t *testing.T, ws *websocket.Conn) {
	if err := ws.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("SetWriteDeadline: %v", err)
	}
	if err := ws.WriteMessage(websocket.TextMessage, []byte("OK OK")); err != nil {
		t.Fatalf("WriteMessage: %v", err)
	}
	if err := ws.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, p, err := ws.ReadMessage()
	if err != nil {
		return
	}
	if !strings.Contains(string(p), "OK") {
		t.Fatalf("message=%s, want %s", p, "OK")
	}
	t.Log("Received=", string(p))
}

func makeWsProto(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
}
