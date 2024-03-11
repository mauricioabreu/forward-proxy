package proxy

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/mauricioabreu/forward-proxy/internal/security"
)

var (
	ErrForbiddenHost = errors.New("host not allowed")
	ErrBannedWord    = errors.New("word not allowed")
)

type Proxy struct {
	// Use a map for fast lookup
	forbiddenHosts map[string]bool
	bannedWords    map[string]bool
}

func New() *Proxy {
	return &Proxy{}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := p.Forward(w, r); err != nil {
		if errors.Is(err, ErrForbiddenHost) {
			http.Error(w, "Access to the requests host is forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Proxy) WithForbiddenHosts(hosts []string) *Proxy {
	p.forbiddenHosts = make(map[string]bool)

	for _, host := range hosts {
		p.forbiddenHosts[host] = true
	}

	return p
}

func (p *Proxy) WithBannedWords(words []string) *Proxy {
	p.bannedWords = make(map[string]bool)

	for _, word := range words {
		p.bannedWords[strings.ToLower(word)] = true
	}

	return p
}

func (p *Proxy) Forward(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodConnect {
		return p.HandleHTTPS(w, r)
	}

	targetURL, err := url.Parse(r.URL.String())
	if err != nil {
		return fmt.Errorf("failed to parse URL: %v", err)
	}

	host, err := extractHost(targetURL)
	if err != nil {
		return fmt.Errorf("failed to host: %v", err)
	}

	if _, forbidden := p.forbiddenHosts[host]; forbidden {
		return ErrForbiddenHost
	}

	req, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	copyRequestHeaders(r, req)

	clientIP := r.RemoteAddr
	if previous, exists := r.Header["X-Forwarded-For"]; exists {
		clientIP = fmt.Sprintf("%s, %s", strings.Join(previous, ", "), clientIP)
	}
	req.Header.Set("X-Forwarded-For", clientIP)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}

	var buf bytes.Buffer
	teeBody := io.TeeReader(resp.Body, &buf)

	body, err := io.ReadAll(teeBody)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if !security.AllowedWord(string(body), p.bannedWords) {
		return ErrBannedWord
	}

	resp.Body = io.NopCloser(&buf)

	copyResponseHeaders(resp, w)

	io.Copy(w, resp.Body)

	return nil
}

func (p *Proxy) HandleHTTPS(w http.ResponseWriter, r *http.Request) error {
	destConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return err
	}
	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return fmt.Errorf("hijacking not supported")
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return fmt.Errorf("hijacking not supported")
	}

	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)

	return nil
}

func transfer(dest io.WriteCloser, source io.ReadCloser) {
	defer dest.Close()
	defer source.Close()
	io.Copy(dest, source)
}

func extractHost(u *url.URL) (string, error) {
	host := u.Host
	hostname, _, err := net.SplitHostPort(host)
	// Host does not have port
	if err != nil {
		return host, nil
	}
	return hostname, nil
}

func isHopHeader(header string) bool {
	hopHeaders := []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"TE",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}

	return slices.Contains(hopHeaders, header)
}

func copyRequestHeaders(from, to *http.Request) {
	for header, values := range from.Header {
		if !isHopHeader(header) {
			for _, value := range values {
				to.Header.Add(header, value)
			}
		}
	}
}

func copyResponseHeaders(from *http.Response, to http.ResponseWriter) {
	for header, values := range from.Header {
		if !isHopHeader(header) {
			for _, value := range values {
				to.Header().Add(header, value)
			}
		}
	}
}
