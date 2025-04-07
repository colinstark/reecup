package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reecup/server"
	"strings"
)

func main() {
	dev := true
	server := server.NewGameServer()

	if dev == true {
		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   "localhost:5173",
		})

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/ws") {
				http.NotFound(w, r)
				return
			}
			proxy.ServeHTTP(w, r)
		})
	} else {
		fs := http.FileServer(http.Dir("./public"))
		http.Handle("/", fs)
	}

	http.HandleFunc("/ws", server.HandleWebSocket)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
