package handler

import "net/http"

func (h ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.R.ServeHTTP(w, r)
}
