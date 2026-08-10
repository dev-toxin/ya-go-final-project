// Package api содержит HTTP-обработчики API планировщика.
package api

import (
	"net/http"
	"os"
)

// Init регистрирует API-маршруты в переданном маршрутизаторе.
func Init(mux *http.ServeMux) {
	config := authConfig{password: os.Getenv("TODO_PASSWORD")}

	mux.HandleFunc("/api/signin", signInHandler(config))
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.Handle("/api/task", auth(config, http.HandlerFunc(taskHandler)))
	mux.Handle("/api/task/done", auth(config, http.HandlerFunc(doneTaskHandler)))
	mux.Handle("/api/tasks", auth(config, http.HandlerFunc(tasksHandler)))
}
