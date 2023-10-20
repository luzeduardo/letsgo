package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// now returns a handler instead of servemux
func (app *application) routes(cfg config) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/sni/view", app.sniView)
	mux.HandleFunc("/sni/create", app.sniCreate)
	// creates a middleware chain
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//return the middleware chain
	return standard.Then(mux)
}
