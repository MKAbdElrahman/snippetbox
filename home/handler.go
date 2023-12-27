package home

import (
	"net/http"
	"snippetbox/error"
	"snippetbox/foundation/logger"
)

type Handler struct {
	errorHandler *error.Handler
}

func NewHandler(logger *logger.Logger) *Handler {
	return &Handler{
		errorHandler: error.NewHandler(logger),
	}
}

func (h *Handler) HandleRenderFullPage(w http.ResponseWriter, r *http.Request) {

	err := HomePage().Render(r.Context(), w)

	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Error rendering full page")
		return
	}

}
