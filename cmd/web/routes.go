package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"snippetbox/foundation/logger"
	"snippetbox/home"
	"snippetbox/httperror"
	"snippetbox/snippet"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(mux *chi.Mux, logger *logger.Logger, db *sql.DB) {

	errorHandelr := httperror.NewHandler(logger)

	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		errorHandelr.NotFound(w, r, "")
	})

	snippetHandler := snippet.NewHandler(logger, db)
	homeHander := home.NewHandler(logger)

	fileServer := http.FileServer(http.Dir("./ui/assets/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", homeHander.HandleRenderFullPage)
	mux.Get("/snippet/view/{id}", snippetHandler.HandleView)
	mux.Get("/snippet/latest", snippetHandler.HandleLatest)

	mux.Get("/snippet/new", snippetHandler.HandleGetNewSnippetForm)

	mux.Post("/snippet/create", snippetHandler.HandleCreate)

	mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		panic(fmt.Errorf("this is a test"))
	})

}
