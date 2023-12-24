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

	ErrorPage(ErrorPageDate{
		Message: err.Message,
		Code:    err.Code,
	}).Render(r.Context(), w)
}

// The convenience methods have an interface suitable for the most common error handling behaviours
// USE HandleHTTPError for better control
// I only log internal server errors because its what i can fix
// change the log flag for logging or debugging during development.

// NotFound responds to the client with a 404 Not Found error.
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	h.HandleHTTPError(w, r, HTTPError{
		Code:    http.StatusNotFound,
		Message: "resource not found",
		Log:     false,
	})
}

// InternalServerError responds to the client with a 500 Internal Server Error.
func (h *Handler) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	h.HandleHTTPError(w, r, HTTPError{
		Code:            http.StatusInternalServerError,
		Message:         http.StatusText(http.StatusInternalServerError),
		UnderlyingError: err,
		Log:             true,
	})
}

// BadRequest responds to the client with a 400 Bad Request error.
func (h *Handler) BadRequest(w http.ResponseWriter, r *http.Request, message string) {
	h.HandleHTTPError(w, r, HTTPError{
		Code:    http.StatusBadRequest,
		Message: message,
		Log:     false,
	})
}

// Unauthorized responds to the client with a 401 Unauthorized error.
func (h *Handler) Unauthorized(w http.ResponseWriter, r *http.Request, message string) {
	h.HandleHTTPError(w, r, HTTPError{
		Code:    http.StatusUnauthorized,
		Message: message,
		Log:     false,
	})
}

// Forbidden responds to the client with a 403 Forbidden error.
func (h *Handler) Forbidden(w http.ResponseWriter, r *http.Request, message string) {
	h.HandleHTTPError(w, r, HTTPError{
		Code:    http.StatusForbidden,
		Message: message,
		Log:     false,
	})
}

// MethodNotAllowed responds to the client with a 405 Method Not Allowed error.
func (h *Handler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	h.HandleHTTPError(w, r, HTTPError{
		Code:    http.StatusMethodNotAllowed,
		Message: "Method Not Allowed",
		Log:     false,
	})
}

