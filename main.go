package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mauricioabreu/forward-proxy/internal/proxy"
)

var (
	forbiddenHostsFile string
	bannedWordsFile    string
)

func init() {
	flag.StringVar(&forbiddenHostsFile, "forbidden-hosts", "", "Forbidden hosts file")
	flag.StringVar(&bannedWordsFile, "banned-words", "", "Banned words file")
}

const (
	serverPort        = "8989"
	readTimeout       = 1 * time.Second
	writeTimeout      = 1 * time.Second
	idleTimeout       = 30 * time.Second
	readHeaderTimeout = 2 * time.Second
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

	if bannedWordsFile != "" {
		bannedWords, err := loadFromFile(bannedWordsFile)
		if err != nil {
			log.Fatalf("Failed to load banned words file: %v", err)
		}

		p.WithBannedWords(bannedWords)
	}

	srv := &http.Server{
		Addr:              ":" + serverPort,
		Handler:           p,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
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
