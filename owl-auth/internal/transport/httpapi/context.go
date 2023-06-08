package httpapi

import (
	"net/http"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}

func ToStdHandler(h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
