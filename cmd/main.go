package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/mkabdelrahman/snippetbox/handler"
	"github.com/mkabdelrahman/snippetbox/view/pages"
)

func home(w http.ResponseWriter, r *http.Request) {
	component := pages.Home()
	err := component.Render(context.Background(), w)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func main() {

	// CONFIG
	const port = 3000
	host := "localhost"
	addr := fmt.Sprintf("%s:%d", host, port)

	// ROUTER
	mux := http.NewServeMux()

	snippetHandler := handler.NewSnippetHandler()

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetHandler.View)
	mux.HandleFunc("GET /snippet/create", snippetHandler.ViewCreateForm)
	mux.HandleFunc("POST /snippet/create", snippetHandler.Create)

	// SERVER
	log.Info("starting server", "host", host, "port", port)
	err := http.ListenAndServe(addr, mux)
	log.Fatal("server error", err)
}
