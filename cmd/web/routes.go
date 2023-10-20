package main

import "net/http"

// now returns a handler instead of servemux
func (app *application) routes(cfg config) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/sni/view", app.sniView)
	mux.HandleFunc("/sni/create", app.sniCreate)
	// pass servemux as the next http.Handler to be executed
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
