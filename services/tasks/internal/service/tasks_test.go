package service

import (
	"tech-ip-sem2/services/tasks/internal/models"
	"tech-ip-sem2/shared/logger"
	"testing"
	"time"
)

func createTestTask(service *TasksService, title, subject string, t *testing.T) models.Task {
	time.Sleep(time.Millisecond * 10)

	task := models.Task{
		Title:       title,
		Description: "Test Description",
		DueDate:     "2026-03-10",
	}

	created, err := service.Create(task, subject)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	return created
}

func TestTasksService(t *testing.T) {
	log := logger.New("test")
	// В тестах используем in-memory хранилище без кэша
	service := NewTasksService(log, nil, nil)

	// Создаем задачу
	created := createTestTask(service, "Test Task", "student", t)

	if created.ID == "" {
		t.Error("Expected task ID to be set")
	}

	if created.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", created.Title)
	}

	// Получаем задачу
	retrieved, err := service.GetByID(created.ID, "student")
	if err != nil {
		t.Errorf("Failed to get task: %v", err)
	}

	if retrieved.ID == "" {
		t.Error("Expected to find task")
	}

	if retrieved.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", retrieved.Title)
	}

	// Получаем все задачи
	tasks, err := service.GetAll("student")
	if err != nil {
		t.Errorf("Failed to get all tasks: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	// Обновляем задачу
	done := true
	updates := models.TaskUpdate{
		Done: &done,
	}

	updated, err := service.Update(created.ID, updates, "student")
	if err != nil {
		t.Errorf("Failed to update task: %v", err)
	}

	if !updated.Done {
		t.Error("Expected task to be done")
	}

	// Удаляем задачу
	deleted, err := service.Delete(created.ID, "student")
	if err != nil {
		t.Errorf("Failed to delete task: %v", err)
	}

	if !deleted {
		t.Error("Expected task to be deleted")
	}

	// Проверяем что удалилось
	retrieved, err = service.GetByID(created.ID, "student")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if retrieved.ID != "" {
		t.Error("Expected empty task after deletion")
	}
}

func TestSearchTasks(t *testing.T) {
	log := logger.New("test")
	service := NewTasksService(log, nil, nil)

	// Создаем несколько задач с разными ID
	tasks := []string{"Go Task", "Docker Task", "K8s Task"}

	for _, title := range tasks {
		_ = createTestTask(service, title, "student", t)
	}

	// Получаем все задачи для проверки
	allTasks, err := service.GetAll("student")
	if err != nil {
		t.Errorf("Failed to get all tasks: %v", err)
	}

	t.Logf("Created %d tasks", len(allTasks))
	for _, task := range allTasks {
		t.Logf("Task: %s (ID: %s)", task.Title, task.ID)
	}

	// Тестируем поиск - в in-memory версии поиск по contains
	results, err := service.SearchByTitle("Go", "student")
	if err != nil {
		t.Errorf("Search failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'Go', got %d", len(results))
		for _, r := range results {
			t.Logf("Found: %s", r.Title)
		}
	} else if results[0].Title != "Go Task" {
		t.Errorf("Expected 'Go Task', got '%s'", results[0].Title)
	}

	// Поиск по "Task" должен найти все
	results, err = service.SearchByTitle("Task", "student")
	if err != nil {
		t.Errorf("Search failed: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results for 'Task', got %d", len(results))
	}
}

// Тест для проверки кэша
func TestCacheWithRealRedis(t *testing.T) {
	t.Skip("Skipping Redis cache test - requires running Redis server")
}
