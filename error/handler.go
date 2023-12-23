package error

import (
	"fmt"
	"net/http"
	"snippetbox/foundation/logger"
)

type HTTPError struct {
	Code            int
	Message         string
	UnderlyingError error
	Log             bool
}

type Handler struct {
	Logger *logger.Logger
}

func NewHandler(logger *logger.Logger) *Handler {
	return &Handler{Logger: logger}
}

// HandleHTTPError handles an HTTPError by logging the error and responding to the client with the specified status code, message, and underlying error.
func (h *Handler) HandleHTTPError(w http.ResponseWriter, r *http.Request, err HTTPError) {
	if err.Log {
		h.Logger.Error(fmt.Sprintf("HTTP Error %d[%s]: %s", err.Code, http.StatusText(err.Code), err.Message), map[string]interface{}{"error": err.UnderlyingError})
	}
	http.Error(w, err.Message, err.Code)
}
