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

	// checks the incmoning request for a session cookie
	//if present , reads the session cookie and retrieves the cooresponding session data from the DB
	// then adds the data to the request context, so it can be used in the handlers
	// also any changes from the handlers are updated in the request context and the middleware updates the DB
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
