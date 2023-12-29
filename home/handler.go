package home

import (
	"net/http"
	"snippetbox/foundation/logger"
	"snippetbox/httperror"
)

type Handler struct {
	errorHandler *httperror.Handler
}

func NewHandler(logger *logger.Logger) *Handler {
	return &Handler{
		errorHandler: httperror.NewHandler(logger),
	}
}

func (h *Handler) HandleRenderFullPage(w http.ResponseWriter, r *http.Request) {

	err := HomePage().Render(r.Context(), w)

	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Error rendering full page")
		return
	}

}
