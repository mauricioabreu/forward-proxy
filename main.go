package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/mauricioabreu/forward-proxy/internal/proxy"
)

var (
	forbiddenHosts     []string
	forbiddenHostsFile string
)

func init() {
	flag.StringVar(&forbiddenHostsFile, "forbidden-hosts", "", "Forbidden hosts file")
}

const (
	serverPort = "8989"
)

func main() {
	flag.Parse()

	p := proxy.New()

	if forbiddenHostsFile != "" {
		forbiddenHosts, err := loadFromFile(forbiddenHostsFile)
		if err != nil {
			log.Fatalf("Failed to load forbidden hosts file: %v", err)
		}

		p.WithForbiddenHosts(forbiddenHosts)
	}

	http.HandleFunc("/", handler(p))

	log.Printf("Running server on %s\n", serverPort)

	if err := http.ListenAndServe(":"+serverPort, nil); err != nil {
		log.Fatalf("Failed to start server; %v", err)
	}
}

func readIntoCollection(reader io.Reader) ([]string, error) {
	var lines []string

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func loadFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	lines, err := readIntoCollection(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %v", err)
	}

	return lines, nil
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
