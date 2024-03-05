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

	assert.NoError(t, proxy.Forward(recorder, req))
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "192.168.0.1", xForwardedFor)
}
