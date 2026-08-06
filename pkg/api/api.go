// Package api содержит HTTP-обработчики API планировщика.
package api

import "net/http"

// Init регистрирует API-маршруты в переданном маршрутизаторе.
func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/task", taskHandler)
	mux.HandleFunc("/api/task/done", doneTaskHandler)
	mux.HandleFunc("/api/tasks", tasksHandler)
}
