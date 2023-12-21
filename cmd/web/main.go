package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	host := "localhost"
	port := 3000

	addr := fmt.Sprintf("%s:%d", host, port)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/assets/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Println("starting server on addr:", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
