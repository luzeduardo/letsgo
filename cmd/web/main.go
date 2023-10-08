package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	addr string
	staticDir string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	//log.New(destination, prefix, additional info)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Llongfile)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/sni/view", sniView)
	mux.HandleFunc("/sni/create", sniCreate)

	infoLog.Printf("Starting server on %s", cfg.addr)
	err := http.ListenAndServe(cfg.addr, mux)
	errorLog.Fatal(err)
}
