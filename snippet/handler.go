package snippet

import (
	"fmt"
	"net/http"
	"snippetbox/error"
	"snippetbox/foundation/logger"
	"strconv"
)

type SnippetHandler struct {
	errorHandler error.Handler
}

func NewHandler(logger *logger.Logger) *SnippetHandler {
	return &SnippetHandler{
		errorHandler: *error.NewHandler(logger),
	}
}

func (h *SnippetHandler) HandleView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		h.errorHandler.HandleHTTPError(w, r, error.HTTPError{
			UnderlyingError: err,
			Message:         "couldn't parse the id query param",
			Code:            http.StatusBadRequest,
			Log:             false,
		})
		return
	}
	if id < 1 {
		h.errorHandler.HandleHTTPError(w, r, error.HTTPError{
			UnderlyingError: err,
			Message:         "resource not found",
			Code:            http.StatusBadRequest,
			Log:             false,
		})
		return
	}

	fmt.Fprintf(w, "Snippet %d", id)
}

func (h *SnippetHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		h.errorHandler.HandleHTTPError(w, r, error.HTTPError{
			UnderlyingError: nil,
			Message:         "only accept get requests",
			Code:            http.StatusBadRequest,
			Log:             false,
		})
		return
	}
	w.Write([]byte("snippet created..."))
}
