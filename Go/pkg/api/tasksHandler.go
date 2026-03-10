package api

import (
	"encoding/json"
	"go_finsl_project/Go/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks *[]db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTasksHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		writeJson(w, map[string]string{"error": ("Wrong HTTP method used")})
		return
	}
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := db.GetTasks(50)

	if err != nil {
		writeJson(w, map[string]string{"error": "error in getting tasks" + err.Error()})
		return
	}
	writeJson(w, TasksResp{
		Tasks: &tasks,
	})
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}
