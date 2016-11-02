package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

const (
	defaultLogFile    = "/tmp/lserver.log"
	defaultListenPort = ":8080"
)

func main() {
	{
		f, err := os.OpenFile(getOrDefault("LOG_FILE", defaultLogFile), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("error:", err)
		}

		log.SetOutput(io.MultiWriter(f, os.Stdout))
	}

	lp := getOrDefault("LISTEN_PORT", defaultListenPort)
	log.Println("Starting listening in " + lp)

	err := http.ListenAndServe(lp, createServer())
	if err != nil {
		log.Fatal("Error ListenAndServe: ", err)
	}
}

func createServer() http.Handler {
	s := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if val := recover(); val != nil {
				log.Printf("[-] Recovering: %+v\nrequest: %+v", val, r)
				debug.PrintStack()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		if r.Method != http.MethodGet {
			log.Println(http.StatusText(http.StatusMethodNotAllowed))
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		log.Printf("Request %s %q remote IP: %q", r.Method, r.URL, r.RemoteAddr)

		gameWS(w, r)
	})

	// route
	http.Handle("/game", s)

	return s
}

func getOrDefault(varName, defaultVal string) string {
	val := os.Getenv(varName)
	if val == "" {
		val = defaultVal
	}
	return val
}
