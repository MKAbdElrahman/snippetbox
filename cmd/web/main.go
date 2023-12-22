package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

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
		}
		log.Fatalf("parsing config: %v", err)
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/assets/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Println("starting server on addr:", cfg.Web.Host)
	err = http.ListenAndServe(cfg.Web.Host, mux)
	if err != nil {
		log.Fatal(err)
	}
}
