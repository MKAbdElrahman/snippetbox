package main

import (
	"context"
	"flag"
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
	addr := flag.String("addr", ":3000", "HTTP network address")
	flag.Parse()

	// ROUTER
	mux := http.NewServeMux()

	// home page
	mux.HandleFunc("GET /{$}", home)

	// static files
	fileServer := http.FileServer(http.Dir("./view/pages/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// snippets
	snippetHandler := handler.NewSnippetHandler()
	mux.HandleFunc("GET /snippet/view/{id}", snippetHandler.View)
	mux.HandleFunc("GET /snippet/create", snippetHandler.ViewCreateForm)
	mux.HandleFunc("POST /snippet/create", snippetHandler.Create)

	// SERVER
	log.Info("starting server", "addr", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal("server error", err)
}
