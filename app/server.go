package app

import (
	"log"
	"net/http"
)

func Run(logger *log.Logger) error {
	mux := http.NewServeMux()
	handler := NewHandler("files", logger)
	staticHandler := http.StripPrefix("/static", http.FileServer(http.Dir("app/static/")))

	mux.Handle("/", handler)
	mux.Handle("/static/", staticHandler)

	server := http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	logger.Printf("Сервер запущен на :8000\n")
	return server.ListenAndServe()
}
