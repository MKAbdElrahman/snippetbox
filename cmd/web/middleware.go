package main

import (
	"fmt"
	"net/http"
)

func SecureHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Content-Security-Policy",
		// 	"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func LogMiddleware(logHandler http.HandlerFunc) func(http.Handler) http.Handler {
	m := func(next http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			logHandler.ServeHTTP(w, r)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(f)
	}
	return m
}

type ErrorHandler interface {
	HandleError(w http.ResponseWriter, r *http.Request, err error)
}

func PanicRecoverMiddleware(errorHandler ErrorHandler) func(http.Handler) http.Handler {
	m := func(next http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connection", "close")
					errorHandler.HandleError(w, r, fmt.Errorf("%s", err))
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(f)
	}
	return m
}
