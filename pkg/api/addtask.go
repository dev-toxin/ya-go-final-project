package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJSON(w, map[string]string{"error": "метод не поддерживается"})
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := db.GetTask(r.URL.Query().Get("id"))
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, task)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := decodeTask(r, &task); err != nil {
		writeJSON(w, map[string]string{"error": fmt.Sprintf("ошибка JSON: %v", err)})
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": fmt.Sprintf("не удалось добавить задачу: %v", err)})
		return
	}
	writeJSON(w, map[string]string{"id": fmt.Sprint(id)})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if err := db.DeleteTask(r.URL.Query().Get("id")); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]string{})
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := decodeTask(r, &task); err != nil {
		writeJSON(w, map[string]string{"error": fmt.Sprintf("ошибка JSON: %v", err)})
		return
	}
	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"})
		return
	}
	if strings.TrimSpace(task.Title) == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]string{})
}

func decodeTask(r *http.Request, task *db.Task) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(task)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateLayout)
	}

	date, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return fmt.Errorf("некорректная дата задачи: %w", err)
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if date.Before(today) {
		if task.Repeat == "" {
			task.Date = now.Format(dateLayout)
		} else {
			task.Date = next
		}
	}
	return nil
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
