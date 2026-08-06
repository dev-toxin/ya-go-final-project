// Package api содержит HTTP-обработчики API планировщика.
package api

import "net/http"

// Init регистрирует API-маршруты в переданном маршрутизаторе.
func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/signin", signInHandler)
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.Handle("/api/task", auth(http.HandlerFunc(taskHandler)))
	mux.Handle("/api/task/done", auth(http.HandlerFunc(doneTaskHandler)))
	mux.Handle("/api/tasks", auth(http.HandlerFunc(tasksHandler)))
}
