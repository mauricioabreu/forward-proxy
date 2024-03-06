package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/mauricioabreu/forward-proxy/proxy"
)

const (
	serverPort = "8989"
)

func main() {
	p := proxy.New()
	http.HandleFunc("/", handler(p))

	log.Printf("Running server on %s\n", serverPort)

	if err := http.ListenAndServe(":"+serverPort, nil); err != nil {
		log.Fatalf("Failed to start server; %v", err)
	}
}

func handler(p *proxy.Proxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := p.Forward(w, r); err != nil {
			if errors.Is(err, proxy.ErrForbiddenHost) {
				http.Error(w, "Access to the requests host is forbidden", http.StatusForbidden)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
