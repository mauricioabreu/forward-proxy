package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

func Forward(w http.ResponseWriter, r *http.Request) error {
	targetURL, err := url.Parse(r.URL.String())
	if err != nil {
		return fmt.Errorf("failed to parse URL: %v", err)
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

	copyResponseHeaders(resp, w)

	io.Copy(w, resp.Body)

	return nil
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
