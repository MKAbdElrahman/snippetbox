package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/mkabdelrahman/snippetbox/central/errorhandler"
	"github.com/mkabdelrahman/snippetbox/handler"
)

func main() {

	// CONFIG
	var config struct {
		addr string
	}

	flag.StringVar(&config.addr, "addr", ":3000", "HTTP network address")
	flag.Parse()

	// LOGGER
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}))

	// ERROR HANDLER
	centralErrorHandler := errorhandler.NewCentralErrorHandler(logger)
	// ROUTER
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

	// SERVER
	logger.Info("starting server", slog.String("addr", config.addr))
	err := http.ListenAndServe(config.addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
