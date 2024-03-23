package main

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func main() {

	const port = 3000
	host := "localhost"
	addr := fmt.Sprintf("%s:%d", host, port)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", home)

	log.Info("starting server", "host", host, "port", port)
	err := http.ListenAndServe(addr, mux)
	log.Fatal("server error", err)
}
