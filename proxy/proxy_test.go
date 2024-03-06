package proxy_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mauricioabreu/forward-proxy/proxy"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	var xForwardedFor string
	remoteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xForwardedFor = r.Header.Get("X-Forwarded-For")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("!P R O X I E D!"))
	}))
	defer remoteServer.Close()

	req, err := http.NewRequest("GET", remoteServer.URL, nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	req.RemoteAddr = "192.168.0.1"

	p := proxy.New()
	assert.NoError(t, p.Forward(recorder, req))
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "192.168.0.1", xForwardedFor)

	body := recorder.Body.String()
	assert.Equal(t, "!P R O X I E D!", body)
}

func TestGetOnDeniedDomain(t *testing.T) {
	remoteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("!P R O X I E D!"))
	}))
	defer remoteServer.Close()

	req, err := http.NewRequest("GET", "http://127.0.0.1", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	p := proxy.New().WithForbiddenHosts([]string{remoteServer.URL})

	assert.NoError(t, p.Forward(recorder, req))
}
