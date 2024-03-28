package errorhandler

import (
	"log/slog"
	"net/http"
)

type CentralErrorHandler struct {
	logger *slog.Logger
}

func NewCentralErrorHandler(logger *slog.Logger) *CentralErrorHandler {
	return &CentralErrorHandler{
		logger: logger,
	}
}
func (h *CentralErrorHandler) HandleInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()
	h.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri))
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

func (h *CentralErrorHandler) HandleBadRequestFromClient(w http.ResponseWriter, r *http.Request, err error) {
	// method := r.Method
	// uri := r.URL.RequestURI()
	// h.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri))
	http.Error(w, "bad request", http.StatusBadRequest)
}
