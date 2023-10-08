package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
}

// holds the app-wide dependencies for the application
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	//log.New(destination, prefix, additional info)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Llongfile)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/sni/view", app.sniView)
	mux.HandleFunc("/sni/create", app.sniCreate)

	infoLog.Printf("Starting server on %s", cfg.addr)
	// by default http server logs error to stdout
	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
