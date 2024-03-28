package main

import (
	"log/slog"
	"net/http"

	"github.com/mkabdelrahman/snippetbox/central/errorhandler"
	"github.com/mkabdelrahman/snippetbox/handler"
)

func buildApplicationRouter(logger *slog.Logger, centralErrorHandler *errorhandler.CentralErrorHandler) http.Handler {
	mux := http.NewServeMux()

	// home page
	homeHandler := handler.NewHomeHandler(logger, centralErrorHandler)
	mux.HandleFunc("GET /{$}", homeHandler.GetHomePage)

	// static files
	fileServer := http.FileServer(http.Dir("./view/pages/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// snippets
	snippetHandler := handler.NewSnippetHandler(logger)
	mux.HandleFunc("GET /snippet/view/{id}", snippetHandler.View)
	mux.HandleFunc("GET /snippet/create", snippetHandler.ViewCreateForm)
	mux.HandleFunc("POST /snippet/create", snippetHandler.Create)

	return mux
}
