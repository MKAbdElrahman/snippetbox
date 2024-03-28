package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/mkabdelrahman/snippetbox/central/errorhandler"
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

	// CENTRAL ERROR HANDLER
	centralErrorHandler := errorhandler.NewCentralErrorHandler(logger)
	// ROUTER
	mux := buildApplicationRouter(logger, centralErrorHandler)

	// SERVER
	logger.Info("starting server", slog.String("addr", config.addr))
	err := http.ListenAndServe(config.addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
