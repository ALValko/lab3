package app

import "net/http"

func InternalServerError(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func NotFoundError(w http.ResponseWriter) {
	http.Error(w, "Page Not Found", http.StatusNotFound)
}
