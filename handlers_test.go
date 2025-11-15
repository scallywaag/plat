package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", GetTasks)
	mux.HandleFunc("GET /tasks/{id}", GetTask)
	mux.HandleFunc("POST /tasks", CreateTask)
	mux.HandleFunc("PUT /tasks/{id}", UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", DeleteTask)

	var rr *httptest.ResponseRecorder
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func resetState() {
	taskMu.Lock()
	tasks = []Task{}
	nextID = 1
	taskMu.Unlock()
}

func TestCRUD(t *testing.T) {
	t.Run("Create Task", func(t *testing.T) {
		resetState()

		payload := `{"title": "Learn Go", "done": false}`
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
		w := executeRequest(req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}

		var task Task
		if err := json.NewDecoder(w.Body).Decode(&task); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if task.ID != 1 || task.Title != "Learn Go" || task.Done {
			t.Errorf("Task not created correctly: %+v", task)
		}

		if loc := w.Header().Get("Location"); loc != "/tasks/1" {
			t.Errorf("Expected Location header /tasks/1, got %q", loc)
		}
	})

	t.Run("Get All Tasks", func(t *testing.T) {
		resetState()
		taskMu.Lock()
		tasks = append(tasks, Task{ID: 1, Title: "Seeded Task", Done: false})
		taskMu.Unlock()

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := executeRequest(req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var tasks []Task
		if err := json.NewDecoder(w.Body).Decode(&tasks); err != nil {
			t.Fatalf("Failed to decode tasks: %v", err)
		}

		if len(tasks) != 1 || tasks[0].Title != "Seeded Task" {
			t.Errorf("Expected 1 task with title 'Seeded Task', got %+v", tasks)
		}
	})

	t.Run("Get Task by ID", func(t *testing.T) {
		resetState()
		taskMu.Lock()
		tasks = append(tasks, Task{ID: 1, Title: "Test Task", Done: false})
		taskMu.Unlock()

		req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
		w := executeRequest(req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var task Task
		if err := json.NewDecoder(w.Body).Decode(&task); err != nil {
			t.Fatalf("Failed to decode task: %v", err)
		}

		if task.ID != 1 || task.Title != "Test Task" {
			t.Errorf("Expected task ID=1, Title='Test Task', got %+v", task)
		}
	})

	t.Run("Update Task", func(t *testing.T) {
		resetState()
		taskMu.Lock()
		tasks = append(tasks, Task{ID: 1, Title: "Old Title", Done: false})
		taskMu.Unlock()

		payload := `{"title": "New Title", "done": true}`
		req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(payload))
		w := executeRequest(req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var task Task
		if err := json.NewDecoder(w.Body).Decode(&task); err != nil {
			t.Fatalf("Failed to decode updated task: %v", err)
		}

		if task.Title != "New Title" || !task.Done {
			t.Errorf("Task not updated correctly: %+v", task)
		}
	})

	t.Run("Delete Task", func(t *testing.T) {
		resetState()
		taskMu.Lock()
		tasks = append(tasks, Task{ID: 1, Title: "Doomed Task", Done: false})
		taskMu.Unlock()

		req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
		w := executeRequest(req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", w.Code)
		}

		req2 := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
		w2 := executeRequest(req2)
		if w2.Code != http.StatusNotFound {
			t.Errorf("Expected 404 after delete, got %d", w2.Code)
		}
	})

	t.Run("Get Non-existent Task", func(t *testing.T) {
		resetState()

		req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
		w := executeRequest(req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}

		var errResp map[string]string
		if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
			t.Fatalf("Failed to decode error: %v", err)
		}
		if errResp["error"] != "Task not found" {
			t.Errorf("Expected error 'Task not found', got %q", errResp["error"])
		}
	})

	t.Run("Invalid ID Format", func(t *testing.T) {
		resetState()

		req := httptest.NewRequest(http.MethodGet, "/tasks/abc", nil)
		w := executeRequest(req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for invalid ID, got %d", w.Code)
		}
	})
}
