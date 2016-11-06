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
	r, err := http.Post(tserver.URL+"/game", "application/json", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusCreated {
		t.Fatalf("Expected: %q we got: %q", http.StatusText(http.StatusCreated), http.StatusText(r.StatusCode))
	}

	mr := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&mr); err != nil {
		t.Fatal(err)
	}

	for _, field := range []string{"playerID", "gameURL"} {
		if v, contains := mr[field]; !contains || v == "" {
			t.Fatal("Not found field:" + field)
		}
	}
}

func TestHealthOK(t *testing.T) {
	r, err := http.Get(tserver.URL + "/health")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Fatalf("Expected: %q we got: %q", http.StatusText(http.StatusOK), http.StatusText(r.StatusCode))
	}

	dec := json.NewDecoder(r.Body)
	var status bool

	if err := dec.Decode(&status); err != nil {
		t.Fatal(err)
	}
	if status == false {
		t.Fatal(status)
	}
}
