package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/mkabdelrahman/snippetbox/central/errorhandler"
	"github.com/mkabdelrahman/snippetbox/view/pages"
)

type HomeHandler struct {
	logger       *slog.Logger
	errorHandler *errorhandler.CentralErrorHandler
}

func NewHomeHandler(logger *slog.Logger, errorHandler *errorhandler.CentralErrorHandler) *HomeHandler {
	return &HomeHandler{
		logger: logger,
	}
}

func (h *HomeHandler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	component := pages.Home()
	err := component.Render(context.Background(), w)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "")
		return
	}
}
