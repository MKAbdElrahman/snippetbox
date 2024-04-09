package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mkabdelrahman/snippetbox/central/errorhandler"
	"github.com/mkabdelrahman/snippetbox/handler"
	"github.com/mkabdelrahman/snippetbox/handler/middleware"

	"github.com/mkabdelrahman/snippetbox/service"
)

func buildApplicationRouter(snippetService *service.SnippetService, logger *slog.Logger, centralErrorHandler *errorhandler.CentralErrorHandler) http.Handler {
	mux := http.NewServeMux()

	// Health endpoint
	mux.HandleFunc("/health", httpHealth())

	// home page
	homeHandler := handler.NewHomeHandler(snippetService, logger, centralErrorHandler)
	mux.HandleFunc("GET /{$}", homeHandler.GetHomePage)

	// static files
	fileServer := http.FileServer(http.Dir("./view/pages/static_assets/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// snippets
	snippetHandler := handler.NewSnippetHandler(snippetService, logger, centralErrorHandler)
	mux.HandleFunc("GET /snippet/view/{id}", snippetHandler.View)

	mux.HandleFunc("GET /snippet/create", snippetHandler.ViewCreateForm)
	mux.HandleFunc("POST /snippet/create", snippetHandler.Create)

	requestLogger := middleware.RequestLogger(logger)
	commonHeaders := middleware.CommonHeaders
	return requestLogger(commonHeaders(mux))
}

func httpHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":"ok"}`)
	}
}
