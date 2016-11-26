package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/nats-io/nats"
	"github.com/yanpozka/checkers/store"
)

const (
	defaultLogFile      = "/tmp/wsserver.log"
	defaultListenPort   = ":8080"
	defaultMemcacheHost = "127.0.0.1:11211"

	natsKey = 1
)

//
var (
	hub = newHub()

	natsClient *nats.Conn
	ms         store.Store
)

const storeCtxKey = "store"

func main() {
	defer natsClient.Close()

	{
		f, err := os.OpenFile(getOrDefault("LOG_FILE", defaultLogFile), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("error:", err)
		}

		log.SetOutput(io.MultiWriter(f, os.Stdout))
	}

	lp := getOrDefault("LISTEN_PORT", defaultListenPort)
	log.Println("Start listening on " + lp)

	if err := http.ListenAndServe(lp, createServer()); err != nil {
		log.Fatal("Error ListenAndServe: ", err)
	}
}

func middleware(inner http.HandlerFunc) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if val := recover(); val != nil {
				log.Printf("[+] Recovering: %+v\nRequest %s %q Remote IP: %q\n", val, r.Method, r.URL, r.RemoteAddr)
				debug.PrintStack()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		log.Printf("Request %s %q Remote IP: %q", r.Method, r.URL, r.RemoteAddr)

		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		inner(w, r)
	})
}

func createServer() http.Handler {
	var err error
	natsClient, err = nats.Connect(getOrDefault("NATS_URL", nats.DefaultURL))
	if err != nil {
		log.Fatal("Error on connect to ", err)
	}
	natsClient.Flush()

	parts := strings.Split(getOrDefault("MEMCACHE_PORT", defaultMemcacheHost), ",")
	ms = store.NewMemcacheStore(parts...)

	mux := http.NewServeMux()

	// route
	mux.Handle("/ws/game/", middleware(gameWS)) // /ws/game/:gameID

	go hub.run()

	return mux
}

func getOrDefault(varName, defaultVal string) string {
	val := os.Getenv(varName)
	if val == "" {
		val = defaultVal
	}
	return val
}
