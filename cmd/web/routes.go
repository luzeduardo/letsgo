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

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/sni/view/:id", dynamic.ThenFunc(app.sniView))
	router.Handler(http.MethodGet, "/sni/create", dynamic.ThenFunc(app.sniCreate))
	router.Handler(http.MethodPost, "/sni/create", dynamic.ThenFunc(app.sniCreatePost))
	// creates a middleware chain
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//return the middleware chain
	return standard.Then(router)
}
