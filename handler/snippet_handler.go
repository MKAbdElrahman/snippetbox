package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mkabdelrahman/snippetbox/central/errorhandler"
	"github.com/mkabdelrahman/snippetbox/model"
	"github.com/mkabdelrahman/snippetbox/service"
	"github.com/mkabdelrahman/snippetbox/view/pages"
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

	param := model.NewSnippetParams{}
	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "server error")
	}
	snippet, err := h.snippetService.Insert(param)
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
	component := pages.Snippet(*snippet)
	err = component.Render(context.Background(), w)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "")
		return
	}
}

func (h *SnippetHandler) ViewCreateForm(w http.ResponseWriter, r *http.Request) {
	component := pages.Form(pages.NewSnippetCreateFormData())
	err := component.Render(context.Background(), w)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "")
		return
	}
}
