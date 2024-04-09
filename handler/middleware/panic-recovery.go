package middleware

import (
	"fmt"
	"net/http"

	"github.com/mkabdelrahman/snippetbox/central/errorhandler"
)

func PanicRecoverer(errorHandle *errorhandler.CentralErrorHandler) func(http.Handler) http.Handler {

	f := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connection", "close")
					errorHandle.HandleInternalServerError(w, r, fmt.Errorf("%s", err), "internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}

	return f
}
