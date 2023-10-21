package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// now returns a handler instead of servemux
func (app *application) routes(cfg config) http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/sni/view/:id", app.sniView)
	router.HandlerFunc(http.MethodGet, "/sni/create", app.sniCreate)
	router.HandlerFunc(http.MethodPost, "/sni/create", app.sniCreatePost)
	// creates a middleware chain
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//return the middleware chain
	return standard.Then(router)
}
