package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/justinas/alice"
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
	mux.HandleFunc("GET /snippet/view/{id}", snippetHandler.HandleViewSingleSnippet)

	mux.HandleFunc("GET /snippet/create", snippetHandler.HandleViewCreateSnippetForm)
	mux.HandleFunc("POST /snippet/create", snippetHandler.HandleCreateSnippetFromForm)

	mux.HandleFunc("POST /api/snippet/create", snippetHandler.HandleCreateSnippetFromJSON)

	logRequests := middleware.RequestLogger(logger)
	setCommonHeaders := middleware.CommonHeaders
	recoverFormPanics := middleware.PanicRecoverer(centralErrorHandler)

	applyOnEveryRequest := alice.New(recoverFormPanics, logRequests, setCommonHeaders)

	return applyOnEveryRequest.Then(mux)
}

func httpHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":"ok"}`)
	}
}
