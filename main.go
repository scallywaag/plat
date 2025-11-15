package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", GetTasks)
	mux.HandleFunc("GET /tasks/{id}", GetTask)
	mux.HandleFunc("POST /tasks", CreateTask)
	mux.HandleFunc("PUT /tasks/{id}", UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", DeleteTask)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
