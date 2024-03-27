package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/mkabdelrahman/snippetbox/view/pages"
)

type HomeHandler struct {
	logger *slog.Logger
}

func NewHomeHandler(logger *slog.Logger) *HomeHandler {
	return &HomeHandler{
		logger: logger,
	}
}

func (h *HomeHandler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	component := pages.Home()
	err := component.Render(context.Background(), w)
	if err != nil {
		h.logger.Error(err.Error(), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
