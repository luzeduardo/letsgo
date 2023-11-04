package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"poc.eduardo-luz.eu/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// creates a mock http handler that writes a status ok 200 as the response body
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// pass the mock http handler to the secure middlware
	// because secureHeaders returns a httpHandler, call ServerHTTP passing the ResponseRecorder and http.Request dummy
	secureHeaders(next).ServeHTTP(rr, req)

	rs := rr.Result()

	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(body), "OK")
}
