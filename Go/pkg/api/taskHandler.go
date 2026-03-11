package api

import (
	"encoding/json"
	"fmt"
	"go_finsl_project/Go/pkg/db"
	"net/http"
	"time"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		writeJson(w, map[string]string{"error": ("Wrong HTTP method used")}, 405)
		return
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": ("JSON serializing error: " + fmt.Sprint(err))}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "No title"}, http.StatusBadRequest)
		return
	}

	now := time.Now().Truncate(24 * time.Hour)
	if task.Date == "" {
		task.Date = now.Format(Format)
	}

	taskTime, err := time.Parse(Format, task.Date)
	if err != nil {
		writeJson(w, map[string]string{"error": "Error in parsing, wrong format: " + err.Error()}, http.StatusBadRequest)
		return
	}

	if taskTime.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(Format)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJson(w, map[string]string{"error": "Error in forming next date: " + err.Error()}, http.StatusBadRequest)
				return
			}
			task.Date = next
		}
	} else {
		if task.Repeat != "" {
			_, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJson(w, map[string]string{"error": "Error in forming next date: " + err.Error()}, http.StatusBadRequest)
				return
			}
		}
	}

	id, err := db.AddTask(task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Error in adding into db" + err.Error()}, http.StatusBadRequest)
		return
	}

	writeJson(w, map[string]string{"id": fmt.Sprint(id)}, http.StatusOK)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	id := query.Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Error in getting task: " + err.Error()}, http.StatusBadRequest)
		return
	}
	writeJson(w, task, http.StatusOK)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "JSON decoding error"}, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeJson(w, map[string]string{"error": "No param id in query"}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "No title"}, http.StatusBadRequest)
		return
	}

	now := time.Now().Truncate(24 * time.Hour)
	taskTime, err := time.Parse(Format, task.Date)
	if err != nil {
		writeJson(w, map[string]string{"error": "Invalid date format"}, http.StatusBadRequest)
		return
	}

	if taskTime.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(Format)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
				return
			}
			task.Date = next
		}
	} else if task.Repeat != "" {
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	writeJson(w, map[string]any{}, http.StatusOK)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	id := query.Get("id")
	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Error in deleting task: " + err.Error()}, http.StatusBadRequest)
		return
	}
	writeJson(w, map[string]any{}, http.StatusOK)
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	id := query.Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Error in getting task: " + err.Error()}, http.StatusBadRequest)
		return
	}
	if task.Repeat == "" {
		err := db.DeleteTask(task.ID)
		if err != nil {
			writeJson(w, map[string]string{"error": "Error in deleting task: " + err.Error()}, http.StatusBadRequest)
			return
		}
	}
	task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": "Error in forming next date: " + err.Error()}, http.StatusBadRequest)
		return
	}
	err = db.UpdateDate(task)
	writeJson(w, map[string]any{}, http.StatusOK)
}
