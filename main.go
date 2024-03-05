package main

import (
	"log"
	"net/http"

	"github.com/mauricioabreu/forward-proxy/proxy"
)

const (
	serverPort = "8989"
)

func main() {
	http.HandleFunc("/", handler)
	log.Printf("Running server on %s\n", serverPort)
	if err := http.ListenAndServe(":"+serverPort, nil); err != nil {
		log.Fatalf("Failed to start server; %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if err := proxy.Forward(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
