package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", GetTasks)
	mux.HandleFunc("GET /tasks/{id}", GetTask)
	mux.HandleFunc("POST /tasks", CreateTask)
	mux.HandleFunc("PUT /tasks/{id}", UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", DeleteTask)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			log.Println("failed to write response:", err)
		}
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			log.Println("failed to write response:", err)
		}
	})

	mux.Handle("/metrics", promhttp.Handler())

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
