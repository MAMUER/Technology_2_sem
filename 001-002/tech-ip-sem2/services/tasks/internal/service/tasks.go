package service

import (
	"sync"
	"time"

	"tech-ip-sem2/services/tasks/internal/models"
)

type TasksService struct {
	tasks map[string]models.Task
	mu    sync.RWMutex
}

func NewTasksService() *TasksService {
	return &TasksService{
		tasks: make(map[string]models.Task),
	}
}

func (s *TasksService) Create(task models.Task) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := generateID()
	task.ID = id
	task.Done = false

	s.tasks[id] = task
	return task
}

func (s *TasksService) GetAll() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *TasksService) GetByID(id string) (models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	return task, exists
}

func (s *TasksService) Update(id string, updates models.TaskUpdate) (models.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return models.Task{}, false
	}

	if updates.Title != nil {
		task.Title = *updates.Title
	}
	if updates.Description != nil {
		task.Description = *updates.Description
	}
	if updates.DueDate != nil {
		task.DueDate = *updates.DueDate
	}
	if updates.Done != nil {
		task.Done = *updates.Done
	}

	s.tasks[id] = task
	return task, true
}

func (s *TasksService) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.tasks[id]
	if !exists {
		return false
	}

	delete(s.tasks, id)
	return true
}

func generateID() string {
	return "t" + time.Now().Format("20060102150405")
}
