package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/mkabdelrahman/snippetbox/central/errorhandler"
	"github.com/mkabdelrahman/snippetbox/service"
	"github.com/mkabdelrahman/snippetbox/view/pages"
)

type HomeHandler struct {
	logger       *slog.Logger
	errorHandler *errorhandler.CentralErrorHandler

	snippetService *service.SnippetService
}

func NewHomeHandler(snippetService *service.SnippetService, logger *slog.Logger, errorHandler *errorhandler.CentralErrorHandler) *HomeHandler {
	return &HomeHandler{
		logger:         logger,
		errorHandler:   errorHandler,
		snippetService: snippetService,
	}
}

func (h *HomeHandler) GetHomePage(w http.ResponseWriter, r *http.Request) {

	latestSnippets, err := h.snippetService.GetLatestSnippets(10)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "")
		return
	}
	pageData := &pages.HomePageData{
		Title:          "Home",
		LatestSnippets: latestSnippets,
	}
	component := pages.Home(pageData)
	err = component.Render(context.Background(), w)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "")
		return
	}
}
