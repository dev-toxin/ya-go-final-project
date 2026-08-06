package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

const tasksLimit = 50

type tasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if date, err := time.Parse("02.01.2006", search); err == nil {
		search = date.Format(dateLayout)
	}

	tasks, err := db.Tasks(tasksLimit, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, tasksResponse{Tasks: tasks})
}
