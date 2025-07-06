package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

var manager = NewTaskManager()

func handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	task := NewTask()
	manager.AddTask(task)

	ctx := context.Background()
	go task.Process(ctx)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func handleTaskOperations(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/task/"):]
	task, ok := manager.GetTask(id)

	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	case http.MethodDelete:
		manager.DeleteTask(id)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func startServer() {

	http.HandleFunc("/task", handleTaskCreate)
	http.HandleFunc("/task/", handleTaskOperations)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Println("ошибка при попытке запуска сервера")
	} else {
		log.Println("Server running on :8080")
	}

}
