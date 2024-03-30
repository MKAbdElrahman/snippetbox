package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mkabdelrahman/snippetbox/central/errorhandler"
	"github.com/mkabdelrahman/snippetbox/model"
	"github.com/mkabdelrahman/snippetbox/service"
)

type SnippetHandler struct {
	logger         *slog.Logger
	errorHandler   *errorhandler.CentralErrorHandler
	snippetService *service.SnippetService
}

func NewSnippetHandler(snippetService *service.SnippetService, logger *slog.Logger, errorHandler *errorhandler.CentralErrorHandler) *SnippetHandler {
	return &SnippetHandler{
		logger:         logger,
		snippetService: snippetService,
	}
}

func (h *SnippetHandler) Create(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusCreated)

	params := model.NewSnippetParams{
		Title:   "Test Snippet",
		Content: "Content...",
		Expires: 5,
	}

	snippet, err := h.snippetService.Insert(params)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "server error")
	}
	fmt.Fprintf(w, "%+v", *snippet)
}

func (h *SnippetHandler) View(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		h.errorHandler.HandleBadRequestFromClient(w, r, err, "couldn't parse id")
		return
	}

	snippet, err := h.snippetService.GetById(id)

	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "server error")
		return
	}

	if snippet == nil || id < 1 {
		h.errorHandler.HandleResourceNotFound(w, r, err, "snippet not found")
		return
	}
	fmt.Fprintf(w, "%+v", *snippet)

}

func (h *SnippetHandler) ViewCreateForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a Form for Creating a New Snippet"))
}
