package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

const (
	defaultLogFile      = "/tmp/apiserver.log"
	defaultListenPort   = ":9090"
	defaultMemcacheHost = "127.0.0.1:11211"
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

	err := http.ListenAndServe(lp, createRouter())
	if err != nil {
		log.Fatal("Error ListenAndServe: ", err)
	}
}

func createRouter() http.Handler {
	router := httprouter.New()

	// panic recover
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		log.Printf("Recovering: %+v\nrequest: %+v", v, r)
		debug.PrintStack()
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// not found handler
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Detected Not Found: %q %+v", r.RequestURI, *r)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	commonh := alice.New(loggerMW, commonMW)

	// routers:
	router.Handler(http.MethodGet, "/health", commonh.ThenFunc(health))
	router.Handler(http.MethodPost, "/game", commonh.ThenFunc(createGame))

	return router
}

func getOrDefault(varName, defaultVal string) string {
	val := os.Getenv(varName)
	if val == "" {
		val = defaultVal
	}
	return val
}
