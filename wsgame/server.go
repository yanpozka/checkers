package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

const (
	defaultLogFile    = "/tmp/wsserver.log"
	defaultListenPort = ":8080"
)

//
var hub = newHub()

func main() {
	{
		f, err := os.OpenFile(getOrDefault("LOG_FILE", defaultLogFile), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("error:", err)
		}

		log.SetOutput(io.MultiWriter(f, os.Stdout))
	}

	lp := getOrDefault("LISTEN_PORT", defaultListenPort)
	log.Println("Start listening in " + lp)

	err := http.ListenAndServe(lp, createServer())
	if err != nil {
		log.Fatal("Error ListenAndServe: ", err)
	}
}

func middleware(inner http.HandlerFunc) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if val := recover(); val != nil {
				log.Printf("[-] Recovering: %+v\nrequest: %+v", val, r)
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
	mux := http.NewServeMux()

	// route
	mux.Handle("/game/", middleware(gameWS))

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
