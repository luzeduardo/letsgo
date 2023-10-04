package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("home!"))
}

func sniView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sni show"))
}

func sniCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sni create"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/sni/view", sniView)
	mux.HandleFunc("/sni/create", sniCreate)
	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
