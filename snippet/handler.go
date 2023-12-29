package snippet

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"snippetbox/foundation/logger"
	"snippetbox/httperror"
	"strconv"

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

func (h *SnippetHandler) HandleView(w http.ResponseWriter, r *http.Request) {

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

	err = ViewSnippet(data, "WithLayout").Render(r.Context(), w)
	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Error rendering snippet view")
		return
	}
}

func (h *SnippetHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.errorHandler.MethodNotAllowed(w, r, "Method not allowed")
	}

	params := NewModelParams{
		Title:          "O snail",
		Content:        "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa",
		DaystToExpires: 7,
	}

	id, err := h.Service.Insert(params)
	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Error creating snippet")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

func (h *SnippetHandler) HandleLatest(w http.ResponseWriter, r *http.Request) {
	snippets, err := h.Service.Latest()
	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Error retrieving latest snippets")
		return
	}

	format := r.Header.Get("Format")

	for _, snippet := range snippets {
		data := NewViewData(snippet)
		err = ViewSnippet(data, format).Render(r.Context(), w)
		if err != nil {
			h.errorHandler.InternalServerError(w, r, err, "Error rendering snippet view")
			return
		}
	}

}

func (h *SnippetHandler) HandleGetNewSnippetForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
}
