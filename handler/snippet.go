package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type SnippetHandler struct {
	logger *slog.Logger
}

func NewSnippetHandler(logger *slog.Logger) *SnippetHandler {
	return &SnippetHandler{
		logger: logger,
	}
}

func (h *SnippetHandler) Create(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}

func (h *SnippetHandler) View(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display Snippet with Id = %d", id)
}

func (h *SnippetHandler) ViewCreateForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a Form for Creating a New Snippet"))
}
