package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testClient(handleFn http.HandlerFunc) (*IasClient, *httptest.Server) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		handleFn.ServeHTTP(w, r)

	}))
	srvUrl, _ := url.Parse(srv.URL)

	client := NewClient(srv.Client(), srvUrl)
	return NewIasClient(client), srv
}

func assertCall[I interface{}](t *testing.T, r *http.Request, expectedPath string, expectedMethod string, expectedBody any) {
	t.Helper()

	var obj I

	if err := json.NewDecoder(r.Body).Decode(&obj); err == nil {
		assert.Equal(t, expectedBody, obj)
	}

	assert.Equal(t, expectedPath, r.URL.Path)
	assert.Equal(t, expectedMethod, r.Method)
}
