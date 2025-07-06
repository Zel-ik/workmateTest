package main

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	Created    TaskStatus = "created"
	Processing TaskStatus = "processing"
	Finished   TaskStatus = "finished"
	Failed     TaskStatus = " failed"
)

type Task struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Status        TaskStatus `json:"status"`
	CreatedTime   time.Time  `json:"created_time"`
	CompletedTime time.Time  `json:"completed_time,omitempty"`
	Duration      string     `json:"duration,omitempty"`
}

func NewTask() *Task {
	task := Task{
		ID:          uuid.New().String(),
		Status:      Created,
		CreatedTime: time.Now(),
	}
	return &task
}

type TaskManager struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*Task),
	}
}

func (tm *TaskManager) AddTask(t *Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tasks[t.ID] = t
}

func (tm *TaskManager) GetTask(id string) (*Task, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	t, ok := tm.tasks[id]
	return t, ok
}

func (tm *TaskManager) DeleteTask(id string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	delete(tm.tasks, id)
}

// пока так, но по сути тут должен происходить весь процессинг
func (t *Task) Process(ctx context.Context) {
	t.Status = Processing
	delay := time.Duration(1+rand.Intn(2)) * time.Minute

	select {
	case <-time.After(delay):
		t.Status = Finished
		t.CompletedTime = time.Now()
		t.Duration = t.CompletedTime.Sub(t.CreatedTime).String()
	case <-ctx.Done():
	}
}
