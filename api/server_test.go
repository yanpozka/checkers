package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var tserver *httptest.Server

func TestMain(m *testing.M) {
	tserver = httptest.NewServer(createRouter())
	defer tserver.Close()

	m.Run()
}

func TestCreateGame(t *testing.T) {
	res, err := http.Post(tserver.URL+"/game", "application/json", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Expected: %q we got: %q", http.StatusText(http.StatusCreated), http.StatusText(res.StatusCode))
	}
}

func TestHealthOK(t *testing.T) {
	res, err := http.Get(tserver.URL + "/health")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected: %q we got: %q", http.StatusText(http.StatusOK), http.StatusText(res.StatusCode))
	}

	dec := json.NewDecoder(res.Body)
	var status bool

	if err := dec.Decode(&status); err != nil {
		t.Fatal(err)
	}
	if status == false {
		t.Fatal(status)
	}
}
