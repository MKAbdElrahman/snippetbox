package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"

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
	logger.Info("starting server", slog.String("addr", config.addr))
	err := http.ListenAndServe(config.addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
