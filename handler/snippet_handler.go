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
		errorHandler:   errorHandler,
	}
}

func (h *SnippetHandler) HandleCreateSnippetFromJSON(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusCreated)

	param := model.NewSnippetParams{}
	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "server error")
		return
	}
	snippet, err := h.snippetService.Insert(param)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "server error")
		return
	}
	fmt.Fprintf(w, "%+v", *snippet)
}

func (h *SnippetHandler) HandleCreateSnippetFromForm(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		h.errorHandler.HandleBadRequestFromClient(w, r, err, "bad request")
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		h.errorHandler.HandleBadRequestFromClient(w, r, err, "bad request")
		return
	}

	param := model.NewSnippetParams{
		Title:   title,
		Content: content,
		Expires: expires,
	}
	snippet, err := h.snippetService.Insert(param)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "server error")
		return
	}
	fmt.Println(snippet.ID)
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", snippet.ID), http.StatusSeeOther)
}

func (h *SnippetHandler) HandleViewSingleSnippet(w http.ResponseWriter, r *http.Request) {

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

func (h *SnippetHandler) HandleViewCreateSnippetForm(w http.ResponseWriter, r *http.Request) {
	component := pages.Form(pages.NewSnippetCreateFormData())
	err := component.Render(context.Background(), w)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "")
		return
	}
}
