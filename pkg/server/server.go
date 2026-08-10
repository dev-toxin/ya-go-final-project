// Package server содержит настройку и запуск HTTP-сервера приложения.
package server

import (
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
)

const defaultPort = "7540"

// Start запускает HTTP-сервер, раздающий статические файлы фронтенда.
func Start() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	mux := http.NewServeMux()
	api.Init(mux)
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	log.Printf("сервер запущен на http://localhost:%s", port)
	return http.ListenAndServe(":"+port, mux)
}
