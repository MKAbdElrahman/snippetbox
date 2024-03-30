package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mkabdelrahman/snippetbox/central/errorhandler"
	"github.com/mkabdelrahman/snippetbox/db"
	"github.com/mkabdelrahman/snippetbox/service"
)

func main() {

	// LOGGER
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// DATABASE
	dsn := os.Getenv("USER_DB_DSN_FOR_SERVER")
	dbConn, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()
	logger.Info("connected to database")
	
	snippetStore := db.NewSnippetStore(dbConn)
	snippetService := service.NewSnippetService(snippetStore)

	// CENTRAL ERROR HANDLER
	centralErrorHandler := errorhandler.NewCentralErrorHandler(logger)
	// ROUTER
	mux := buildApplicationRouter(snippetService, logger, centralErrorHandler)

	// SERVER
	addr := os.Getenv("SERVER_ADDR")
	logger.Info("starting server", slog.String("addr", addr))
	err = http.ListenAndServe(addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
