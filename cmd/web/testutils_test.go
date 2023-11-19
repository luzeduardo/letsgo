package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"poc.eduardo-luz.eu/internal/models/mocks"
)

// returnsan instance of the app struct with mocked deps
func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		errorLog:       log.New(io.Discard, "", 0),
		infoLog:        log.New(io.Discard, "", 0),
		snippets:       &mocks.SnippetModel{},
		users:          &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	//disable redirect-follow for the test server and client by setting a custom check redirect function
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		//forces the client to immediately return the received response by always returning a ErrUseLastResponse
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// implement a get method on the custom testServer.
// it mkaes a GET request to a given url using the testclient and returns the response status code
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	route := ts.URL + urlPath
	rs, err := ts.Client().Get(route)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
