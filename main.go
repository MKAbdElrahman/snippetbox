package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a snippet."))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a form to create a new snippet."))
}

func main() {

	host := "localhost"
	port := 3000

	addr := fmt.Sprintf("%s:%d", host, port)

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Println("starting server on addr:", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
