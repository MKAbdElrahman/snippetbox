package snippet

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"snippetbox/foundation/logger"
	"snippetbox/httperror"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
)

type SnippetHandler struct {
	errorHandler *httperror.Handler
	Service      *Service
}

func NewHandler(logger *logger.Logger, db *sql.DB) *SnippetHandler {
	return &SnippetHandler{
		errorHandler: httperror.NewHandler(logger),
		Service:      NewService(db),
	}
}

func (h *SnippetHandler) ViewSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		h.errorHandler.BadRequest(w, r, "Invalid ID parameter")
		return
	}
	if id < 1 {
		h.errorHandler.NotFound(w, r, "Snippet not found")
		return
	}

	snippet, err := h.Service.Get(id)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			h.errorHandler.NotFound(w, r, "Snippet not found")
			return
		} else {
			h.errorHandler.InternalServerError(w, r, err, "Error retrieving snippet")
			return
		}
	}
	data := NewViewData(snippet)

	err = ViewSnippet(data).Render(r.Context(), w)
	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Error rendering snippet view")
		return
	}
}

func (h *SnippetHandler) CreateSnippet(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {

		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	expirationDate, err := time.Parse("2006-01-02", r.FormValue("expiresDate"))
	if err != nil {
		http.Error(w, "Error parsing expiration date", http.StatusBadRequest)
		return
	}

	expirationTime, err := time.Parse("15:04", r.FormValue("expiresTime"))
	if err != nil {
		http.Error(w, "Error parsing expiration time", http.StatusBadRequest)
		return
	}

	snippetData := NewModelParams{
		Title:       r.FormValue("title"),
		Content:     r.FormValue("content"),
		ExpiresDate: expirationDate,
		ExpiresTime: expirationTime,
	}

	id, err := h.Service.Insert(snippetData)
	if err != nil {
		fmt.Println(snippetData)
		h.errorHandler.InternalServerError(w, r, err, "Error creating snippet")
		return
	}

	m, err := h.Service.Get(id)
	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Error getting snippet")
		return
	}

	err = ViewSnippet(NewViewData(m)).Render(r.Context(), w)
	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Error rendering snippet view")
		return
	}
}

func (h *SnippetHandler) ListLatestSnippets(w http.ResponseWriter, r *http.Request) {
	snippets, err := h.Service.Latest()
	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Error retrieving latest snippets")
		return
	}

	for _, snippet := range snippets {
		data := NewViewData(snippet)
		err = ViewSnippet(data).Render(r.Context(), w)
		if err != nil {
			h.errorHandler.InternalServerError(w, r, err, "Error rendering snippet view")
			return
		}
	}

}

func (h *SnippetHandler) GetNewSnippetForm(w http.ResponseWriter, r *http.Request) {
	err := CreateSnippetForm().Render(r.Context(), w)

	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Error rendering full page")
		return
	}
}

func (h *SnippetHandler) GetSearchSnippetForm(w http.ResponseWriter, r *http.Request) {
	err := SearchSnippetForm().Render(r.Context(), w)

	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Error rendering full page")
		return
	}
}

func (h *SnippetHandler) DeleteSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the snippet ID from the URL parameter.
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.errorHandler.BadRequest(w, r, "Invalid ID parameter")
		return
	}

	// Check if the ID is valid.
	if id < 1 {
		h.errorHandler.NotFound(w, r, "Snippet not found")
		return
	}

	// Delete the snippet.
	err = h.Service.Delete(id)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			h.errorHandler.NotFound(w, r, "Snippet not found")
			return
		} else {
			h.errorHandler.InternalServerError(w, r, err, "Error deleting snippet")
			return
		}
	}

}

func (h *SnippetHandler) ValidateSnippetTitle(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	if strings.TrimSpace(title) == "" {
		fmt.Fprint(w, "This field cannot be blank")
	}

	if utf8.RuneCountInString(title) > 100 {
		fmt.Fprint(w, "This field cannot be more than 100 characters long")
	}
}
