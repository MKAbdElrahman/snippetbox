package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}
func main() {

	host := "localhost"
	port := 3000

	addr := fmt.Sprintf("%s:%d", host, port)

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)

	log.Println("starting server on addr:", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
