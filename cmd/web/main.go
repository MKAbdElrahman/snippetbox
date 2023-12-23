package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"snippetbox/foundation/logger"
	"snippetbox/home"
	"snippetbox/snippet"

	"github.com/ardanlabs/conf/v3"
)

func main() {

	cfg := struct {
		conf.Version
		Web struct {
			Host string `conf:"default:0.0.0.0:3000"`
		}
	}{
		Version: conf.Version{
			Build: "v1.0.0",
			Desc:  "Snippetbox",
		},
	}

	const prefix = "Snippetbox"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return
		} else {
			fmt.Printf("parsing config: %v", err)
			os.Exit(1)
		}
	}

	logger := logger.NewLogger(
		logger.WithStderr(os.Stderr),
		logger.WithStdout(os.Stdout),
		logger.WithDate(true, logger.ErrorLevel, logger.InfoLevel),
		logger.WithTime(true, logger.InfoLevel, logger.ErrorLevel),
		logger.WithLineSource(true, logger.ErrorLevel, logger.InfoLevel),
	)

	mux := http.NewServeMux()
	RegisterRoutes(mux, logger)

	server := http.Server{
		Handler: mux,
		Addr:    cfg.Web.Host,
	}

	logger.Info("starting server", map[string]any{"Host": cfg.Web.Host})
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("server is not listening")
	}
}

func RegisterRoutes(mux *http.ServeMux, logger *logger.Logger) {

	snippetHandler := snippet.NewHandler(logger)
	homeHander := home.NewHandler(logger)

	fileServer := http.FileServer(http.Dir("./ui/assets/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", homeHander.HandleRenderFullPage)
	mux.HandleFunc("/snippet/view", snippetHandler.HandleView)
	mux.HandleFunc("/snippet/create", snippetHandler.HandleCreate)
}
