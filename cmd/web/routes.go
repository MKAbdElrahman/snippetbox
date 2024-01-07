package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"snippetbox/foundation/logger"
	"snippetbox/home"
	"snippetbox/httperror"
	"snippetbox/snippet"
	"snippetbox/user"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(mux *chi.Mux, logger *logger.Logger, sessionManager *scs.SessionManager, db *sql.DB) {

	errorHandler := httperror.NewHandler(logger)

	// Serve static files
	fileServer := http.FileServer(http.Dir("./ui/assets/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Handle 404 Not Found
	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		errorHandler.NotFound(w, r, "")
	})

	// Home routes
	homeHandler := home.NewHandler(logger)

	// User Routes
	userHandler := user.NewHandler(logger, sessionManager, db)

	mux.Route("/", func(r chi.Router) {
		r.Use(sessionManager.LoadAndSave)
		r.Use(noSurf)
		r.Use(RequireAuthentication(sessionManager))

		r.Get("/", homeHandler.HandleRenderFullPage)
	})

	mux.Route("/user", func(r chi.Router) {

		r.Use(sessionManager.LoadAndSave)
		r.Use(noSurf)

		r.Get("/signup", userHandler.GetUserSignUpForm)
		r.Post("/signup", userHandler.HandleSignUpUser)

		r.Get("/login", userHandler.GetUserLoginForm)
		r.Post("/login", userHandler.HandleLoginUser)

		r.Post("/logout", userHandler.HandleLogoutUser)

	})

	// Snippet routes
	snippetHandler := snippet.NewHandler(logger, sessionManager, db)
	mux.Route("/snippets", func(r chi.Router) {

		r.Use(sessionManager.LoadAndSave)
		r.Use(noSurf)
		r.Use(RequireAuthentication(sessionManager))

		r.Get("/", snippetHandler.ListLatestSnippets)
		r.Get("/{id}", snippetHandler.ViewSnippet)
		r.Delete("/{id}", snippetHandler.DeleteSnippet)

		r.Get("/form/create", snippetHandler.GetNewSnippetForm)
		r.Get("/form/search", snippetHandler.GetSearchSnippetForm)

		r.Post("/", snippetHandler.CreateSnippet)

		r.Post("/title/validate", snippetHandler.ValidateSnippetTitle)

	})

	// Test route
	mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		panic(fmt.Errorf("this is a test"))
	})
}
