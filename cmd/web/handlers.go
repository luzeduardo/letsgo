package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {

if r.URL.Path != "/" {
http.NotFound(w, r)
return
}
w.Write([]byte("home!"))
}

func sniView(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Cache-Control", "public, max-age=31536000")
w.Header()["Date"] = nil

id, err := strconv.Atoi(r.URL.Query().Get("id"))
if err != nil || id < 1 {
http.NotFound(w, r)
return
}
fmt.Fprintf(w, "Display item with ID %d", id)
}

func sniCreate(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
w.Header().Set("Allow", "POST")
http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
return
}
w.Header().Set("Content-Type", "application/json")
w.Write([]byte(`{"name":"JOhn Doe"`))
}
