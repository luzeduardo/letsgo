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
	// Unprotected application routes using the dynamic middleware chain
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)
	// protected authenticated application routes that includes requiredAuthentication middleware
	router.Handler(http.MethodGet, "/sni/view/:id", protected.ThenFunc(app.sniView))
	router.Handler(http.MethodPost, "/sni/create", protected.ThenFunc(app.sniCreatePost))
	router.Handler(http.MethodGet, "/sni/create", protected.ThenFunc(app.sniCreate))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// creates a middleware chain
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	//return the middleware chain
	return standard.Then(router)
}
