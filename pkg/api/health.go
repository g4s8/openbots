package api

import "net/http"

type health struct{}

func (h *health) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
