package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"snippetbox/foundation/logger"
	"snippetbox/home"
	"snippetbox/snippet"
)

func RegisterRoutes(mux *http.ServeMux, logger *logger.Logger, db *sql.DB) {

	snippetHandler := snippet.NewHandler(logger, db)
	homeHander := home.NewHandler(logger)

	fileServer := http.FileServer(http.Dir("./ui/assets/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", homeHander.HandleRenderFullPage)
	mux.HandleFunc("/snippet/view", snippetHandler.HandleView)
	mux.HandleFunc("/snippet/latest", snippetHandler.HandleLatest)
	mux.HandleFunc("/snippet/create", snippetHandler.HandleCreate)

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		panic(fmt.Errorf("this is a test"))
	})

}
