package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"snippetbox/foundation/logger"
	"snippetbox/home"
	"snippetbox/snippet"

	"github.com/ardanlabs/conf/v3"
	"github.com/justinas/alice"

	_ "github.com/go-sql-driver/mysql"
)

func HandleLog(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stdout, r.Method+" "+r.URL.Path+"\n")
}

type MyErrorHandler struct {
}

func (h *MyErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "internal server error", http.StatusInternalServerError)
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
		logger.WithLineSource(true, logger.ErrorLevel, logger.InfoLevel),
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

	mux := http.NewServeMux()
	RegisterRoutes(mux, logger, db)

	middleware := alice.New(PanicRecoverMiddleware(&MyErrorHandler{}), LogMiddleware(HandleLog), SecureHeadersMiddleware)
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

func RegisterRoutes(mux *http.ServeMux, logger *logger.Logger, db *sql.DB) {

	snippetHandler := snippet.NewHandler(logger, db)
	homeHander := home.NewHandler(logger)

	fileServer := http.FileServer(http.Dir("./ui/assets/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", homeHander.HandleRenderFullPage)
	mux.HandleFunc("/snippet/view", snippetHandler.HandleView)
	mux.HandleFunc("/snippet/latest", snippetHandler.HandleLatest)
	mux.HandleFunc("/snippet/create", snippetHandler.HandleCreate)

	mux.HandleFunc("/test", HandleTest)

}

func HandleTest(w http.ResponseWriter, r *http.Request) {
	panic(fmt.Errorf("this is a test"))
}
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
