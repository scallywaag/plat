package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var (
	tasks  = []Task{}
	taskMu = sync.RWMutex{}
	nextID = 1
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	taskMu.RLock()
	defer taskMu.RUnlock()

	jsonResponse(w, http.StatusOK, tasks)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		jsonError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	taskMu.RLock()
	defer taskMu.RUnlock()

	for _, t := range tasks {
		if t.ID == id {
			jsonResponse(w, http.StatusOK, t)
			return
		}
	}
	jsonError(w, http.StatusNotFound, "Task not found")
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if task.Title == "" {
		jsonError(w, http.StatusBadRequest, "Title is required")
		return
	}

	taskMu.Lock()
	task.ID = nextID
	nextID++
	tasks = append(tasks, task)
	taskMu.Unlock()

	w.Header().Set("Location", "/tasks/"+strconv.Itoa(task.ID))
	jsonResponse(w, http.StatusCreated, task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		jsonError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var updates Task
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	taskMu.Lock()
	defer taskMu.Unlock()

	for i, t := range tasks {
		if t.ID == id {
			if updates.Title != "" {
				tasks[i].Title = updates.Title
			}
			tasks[i].Done = updates.Done
			jsonResponse(w, http.StatusOK, tasks[i])
			return
		}
	}
	jsonError(w, http.StatusNotFound, "Task not found")
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		jsonError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	taskMu.Lock()
	defer taskMu.Unlock()

	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	jsonError(w, http.StatusNotFound, "Task not found")
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func jsonError(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}
