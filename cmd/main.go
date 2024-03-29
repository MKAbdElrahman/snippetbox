package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mkabdelrahman/snippetbox/central/errorhandler"
)

func main() {

	// LOGGER
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// DATABASE
	dsn := os.Getenv("DB_DSN_FOR_SERVER")
	db, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	logger.Info("connected to database")

	// CENTRAL ERROR HANDLER
	centralErrorHandler := errorhandler.NewCentralErrorHandler(logger)
	// ROUTER
	mux := buildApplicationRouter(logger, centralErrorHandler)

	// SERVER
	addr := os.Getenv("SERVER_ADDR")
	logger.Info("starting server", slog.String("addr", addr))
	err = http.ListenAndServe(addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
