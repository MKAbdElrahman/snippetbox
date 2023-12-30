package httperror

import (
	"fmt"
	"net/http"
)

// HTTPError represents an HTTP error.
type HTTPError struct {
	Code            int
	UserMessage     string
	UnderlyingError error
	Log             bool
}

// Error implements the error interface for HTTPError.
func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP Error %d: %s", e.Code, e.UserMessage)
}

// WithMessage sets a user-friendly message for the error.
func (e HTTPError) WithMessage(userMessage string) HTTPError {
	e.UserMessage = userMessage
	return e
}

// NewHTTPError creates a new HTTPError with the given status code, optional user message, and optional underlying error.
func NewHTTPError(code int, userMessage string, underlyingError ...error) HTTPError {
	err := HTTPError{
		Code:        code,
		UserMessage: userMessage,
		Log:         true,
	}

	if len(underlyingError) > 0 {
		err.UnderlyingError = underlyingError[0]
	}

	return err
}

// Logger is an interface for logging errors.
type Logger interface {
	Error(msg string, fields ...map[string]interface{})
}

// Handler is responsible for handling HTTP errors.
type Handler struct {
	Logger Logger
}

// NewHandler creates a new instance of Handler.
func NewHandler(logger Logger) *Handler {
	return &Handler{Logger: logger}
}

// HandleError handles a generic error by checking if it can be cast to an HTTPError.
// If not, it defaults to an internal server error.
func (h *Handler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if httpErr, ok := err.(HTTPError); ok {
		h.HandleHTTPError(w, r, httpErr)
	} else {
		h.InternalServerError(w, r, err, "")
	}
}

// HandleHTTPError handles an HTTPError by logging the error and responding to the client with the specified status code, user message, and underlying error.
func (h *Handler) HandleHTTPError(w http.ResponseWriter, r *http.Request, err HTTPError) {
	if err.Log {
		h.Logger.Error(
			fmt.Sprintf("HTTPErr%d", err.Code), map[string]interface{}{"error": err.UnderlyingError},
		)
	}

	ErrorPage(ErrorPageData{
		Message: err.UserMessage,
		Code:    err.Code,
	}).Render(r.Context(), w)
}

// convenienceError is a helper function for handling common HTTP errors.
func (h *Handler) convenienceError(w http.ResponseWriter, r *http.Request, code int, userMessage string) {
	err := HTTPError{
		Code:            code,
		UserMessage:     userMessage,
		UnderlyingError: nil,
		Log:             false,
	}

	h.HandleHTTPError(w, r, err)
}

// NotFound responds to the client with a 404 Not Found error.
func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request, userMessage string) {
	h.convenienceError(w, r, http.StatusNotFound, userMessage)
}

// InternalServerError responds to the client with a 500 Internal Server Error.
func (h *Handler) InternalServerError(w http.ResponseWriter, r *http.Request, err error, userMessage string) {
	httpErr := HTTPError{
		Code:            http.StatusInternalServerError,
		UserMessage:     userMessage,
		UnderlyingError: err,
		Log:             true,
	}

	h.HandleHTTPError(w, r, httpErr)
}

// BadRequest responds to the client with a 400 Bad Request error.
func (h *Handler) BadRequest(w http.ResponseWriter, r *http.Request, userMessage string) {
	h.convenienceError(w, r, http.StatusBadRequest, userMessage)
}

// Unauthorized responds to the client with a 401 Unauthorized error.
func (h *Handler) Unauthorized(w http.ResponseWriter, r *http.Request, userMessage string) {
	h.convenienceError(w, r, http.StatusUnauthorized, userMessage)
}

// Forbidden responds to the client with a 403 Forbidden error.
func (h *Handler) Forbidden(w http.ResponseWriter, r *http.Request, userMessage string) {
	h.convenienceError(w, r, http.StatusForbidden, userMessage)
}

// MethodNotAllowed responds to the client with a 405 Method Not Allowed error.
func (h *Handler) MethodNotAllowed(w http.ResponseWriter, r *http.Request, userMessage string) {
	h.convenienceError(w, r, http.StatusMethodNotAllowed, userMessage)
}
