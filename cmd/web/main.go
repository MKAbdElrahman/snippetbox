package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"snippetbox/foundation/logger"
	"snippetbox/httperror"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/ardanlabs/conf/v3"
	"github.com/justinas/alice"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

func HandleLog(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stdout, r.Method+" "+r.URL.Path+"\n")
}

func main() {

	cfg := struct {
		conf.Version
		Web struct {
			Host string `conf:"default:0.0.0.0:3000"`
		}

		Mysql struct {
			Dsn string `conf:"default:user:password@/snippetbox?parseTime=true"`
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
	)

	fmt.Println(cfg.Mysql.Dsn)
	db, err := openDB(cfg.Mysql.Dsn)
	if err != nil {
		logger.Error("failed to open db", map[string]any{
			"error": err.Error(),
		})
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("connected to mysql", map[string]any{})

	// Session Manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	mux := chi.NewRouter()
	RegisterRoutes(mux, logger, sessionManager, db)

	middleware := alice.New(
		PanicRecoverMiddleware(httperror.NewHandler(logger)),
		LogMiddleware(HandleLog),
		SecureHeadersMiddleware)

	server := http.Server{
		Handler: middleware.Then(mux),
		Addr:    cfg.Web.Host,
	}

	logger.Info("starting server", map[string]any{"Host": cfg.Web.Host})
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("server is not listening")
	}
}
