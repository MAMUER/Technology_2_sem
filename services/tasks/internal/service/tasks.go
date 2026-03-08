package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"tech-ip-sem2/services/tasks/internal/cache"
	"tech-ip-sem2/services/tasks/internal/models"
	"tech-ip-sem2/services/tasks/internal/repository"
	"tech-ip-sem2/shared/logger"
	"tech-ip-sem2/shared/sanitize"
)

type TasksService struct {
	repo        repository.TaskRepository
	cache       *cache.RedisCache
	log         *logger.Logger
	memoryTasks map[string]models.Task
	mu          sync.RWMutex
	useDatabase bool
	counter     int64
}

func NewTasksService(log *logger.Logger, repo repository.TaskRepository, cache *cache.RedisCache) *TasksService {
	log.Info("Tasks service initialized",
		zap.Bool("use_database", repo != nil),
		zap.Bool("cache_enabled", cache != nil && cache.IsEnabled()),
	)

	return &TasksService{
		repo:        repo,
		cache:       cache,
		log:         log,
		memoryTasks: make(map[string]models.Task),
		useDatabase: repo != nil,
		counter:     0,
	}
}

// GetByID с поддержкой кэша (cache-aside)
func (s *TasksService) GetByID(id string, subject string) (models.Task, error) {
	ctx := context.Background()

	// 1. Пытаемся получить из кэша
	if s.cache != nil && s.cache.IsEnabled() {
		cachedTask, err := s.cache.GetTask(ctx, id)
		if err != nil {
			// Ошибка Redis - логируем, но продолжаем работу с БД
			s.log.Warn("Cache read error, falling back to database",
				zap.Error(err),
				zap.String("task_id", id),
			)
		} else if cachedTask != nil {
			// Cache HIT
			if cachedTask.Subject == subject {
				s.log.Debug("Cache hit for task",
					zap.String("task_id", id),
					zap.String("subject", subject),
				)
				return *cachedTask, nil
			}
			// Задача принадлежит другому пользователю - игнорируем кэш
			s.log.Debug("Cache hit but wrong subject",
				zap.String("task_id", id),
				zap.String("cached_subject", cachedTask.Subject),
				zap.String("request_subject", subject),
			)
		} else {
			s.log.Debug("Cache miss for task", zap.String("task_id", id))
		}
	}

	// 2. Cache MISS или ошибка - идем в БД/память
	var task models.Task
	var err error

	if s.useDatabase {
		task, err = s.repo.GetByID(id, subject)
	} else {
		s.mu.RLock()
		task, err = s.getByIDMemory(id, subject)
		s.mu.RUnlock()
	}

	if err != nil {
		return models.Task{}, err
	}

	if task.ID == "" {
		return models.Task{}, nil
	}

	// 3. Сохраняем в кэш
	if s.cache != nil && s.cache.IsEnabled() {
		go func() {
			if err := s.cache.SetTask(ctx, &task); err != nil {
				s.log.Warn("Failed to cache task",
					zap.Error(err),
					zap.String("task_id", task.ID),
				)
			}
		}()
	}

	return task, nil
}

// Вспомогательный метод для получения из памяти
func (s *TasksService) getByIDMemory(id string, subject string) (models.Task, error) {
	task, exists := s.memoryTasks[id]
	if !exists || task.Subject != subject {
		return models.Task{}, nil
	}
	return task, nil
}

// GetAll с поддержкой кэша
func (s *TasksService) GetAll(subject string) ([]models.Task, error) {
	ctx := context.Background()

	// 1. Пытаемся получить список из кэша
	if s.cache != nil && s.cache.IsEnabled() {
		cachedTasks, err := s.cache.GetTaskList(ctx, subject)
		if err != nil {
			s.log.Warn("Cache read error for list, falling back to database",
				zap.Error(err),
				zap.String("subject", subject),
			)
		} else if cachedTasks != nil {
			s.log.Debug("Cache hit for task list",
				zap.String("subject", subject),
				zap.Int("count", len(cachedTasks)),
			)
			return cachedTasks, nil
		} else {
			s.log.Debug("Cache miss for task list", zap.String("subject", subject))
		}
	}

	// 2. Cache MISS - идем в БД/память
	var tasks []models.Task
	var err error

	if s.useDatabase {
		tasks, err = s.repo.GetAll(subject)
	} else {
		s.mu.RLock()
		tasks = s.getAllMemory(subject)
		s.mu.RUnlock()
	}

	if err != nil {
		return nil, err
	}

	// 3. Сохраняем в кэш
	if s.cache != nil && s.cache.IsEnabled() && len(tasks) > 0 {
		go func() {
			if err := s.cache.SetTaskList(ctx, subject, tasks); err != nil {
				s.log.Warn("Failed to cache task list",
					zap.Error(err),
					zap.String("subject", subject),
				)
			}
		}()
	}

	s.log.Debug("Tasks listed", zap.Int("count", len(tasks)), zap.String("subject", subject))
	return tasks, nil
}

// Вспомогательный метод для получения всех задач из памяти
func (s *TasksService) getAllMemory(subject string) []models.Task {
	tasks := make([]models.Task, 0, len(s.memoryTasks))
	for _, task := range s.memoryTasks {
		if task.Subject == subject {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// Create
func (s *TasksService) Create(task models.Task, subject string) (models.Task, error) {
	task.Sanitize()

	var created models.Task
	var err error

	if s.useDatabase {
		task.ID = generateUUID()
		created, err = s.repo.Create(task, subject)
	} else {
		s.mu.Lock()
		created = s.createMemory(task, subject)
		s.mu.Unlock()
	}

	if err != nil {
		return models.Task{}, err
	}

	if s.cache != nil && s.cache.IsEnabled() {
		go func() {
			if err := s.cache.DeleteTaskList(context.Background(), subject); err != nil {
				s.log.Warn("Failed to invalidate task list cache",
					zap.Error(err),
					zap.String("subject", subject),
				)
			}
		}()
	}

	s.log.Info("Task created",
		zap.String("task_id", created.ID),
		zap.String("title", created.Title),
		zap.String("subject", subject),
	)

	return created, nil
}

// Вспомогательный метод для создания в памяти
func (s *TasksService) createMemory(task models.Task, subject string) models.Task {
	s.counter++
	task.ID = s.generateID()
	task.Done = false
	task.Subject = subject
	s.memoryTasks[task.ID] = task
	return task
}

// Update
func (s *TasksService) Update(id string, updates models.TaskUpdate, subject string) (models.Task, error) {
	if updates.Description != nil {
		sanitized, err := sanitize.ValidateAndSanitizeDescription(*updates.Description)
		if err != nil {
			return models.Task{}, err
		}
		updates.Description = &sanitized
	}
	if updates.Title != nil {
		*updates.Title = sanitize.SanitizeText(*updates.Title)
	}

	var updated models.Task
	var err error

	if s.useDatabase {
		updated, err = s.repo.Update(id, updates, subject)
	} else {
		s.mu.Lock()
		updated, err = s.updateMemory(id, updates, subject)
		s.mu.Unlock()
	}

	if err != nil {
		return models.Task{}, err
	}

	if updated.ID == "" {
		return models.Task{}, nil
	}

	if s.cache != nil && s.cache.IsEnabled() {
		go func() {
			ctx := context.Background()
			// Удаляем конкретную задачу из кэша
			if err := s.cache.DeleteTask(ctx, id); err != nil {
				s.log.Warn("Failed to invalidate task cache",
					zap.Error(err),
					zap.String("task_id", id),
				)
			}
			// Удаляем список (так как данные изменились)
			if err := s.cache.DeleteTaskList(ctx, subject); err != nil {
				s.log.Warn("Failed to invalidate task list cache",
					zap.Error(err),
					zap.String("subject", subject),
				)
			}
		}()
	}

	s.log.Info("Task updated", zap.String("task_id", id))
	return updated, nil
}

// Вспомогательный метод для обновления в памяти
func (s *TasksService) updateMemory(id string, updates models.TaskUpdate, subject string) (models.Task, error) {
	task, exists := s.memoryTasks[id]
	if !exists || task.Subject != subject {
		return models.Task{}, nil
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

	s.memoryTasks[id] = task
	return task, nil
}

// Delete
func (s *TasksService) Delete(id string, subject string) (bool, error) {
	var deleted bool
	var err error

	if s.useDatabase {
		deleted, err = s.repo.Delete(id, subject)
	} else {
		s.mu.Lock()
		deleted = s.deleteMemory(id, subject)
		s.mu.Unlock()
	}

	if err != nil {
		return false, err
	}

	if !deleted {
		return false, nil
	}

	if s.cache != nil && s.cache.IsEnabled() {
		go func() {
			ctx := context.Background()
			// Удаляем задачу из кэша
			if err := s.cache.DeleteTask(ctx, id); err != nil {
				s.log.Warn("Failed to invalidate task cache on delete",
					zap.Error(err),
					zap.String("task_id", id),
				)
			}
			// Удаляем список
			if err := s.cache.DeleteTaskList(ctx, subject); err != nil {
				s.log.Warn("Failed to invalidate task list cache on delete",
					zap.Error(err),
					zap.String("subject", subject),
				)
			}
		}()
	}

	s.log.Info("Task deleted", zap.String("task_id", id))
	return true, nil
}

// Вспомогательный метод для удаления из памяти
func (s *TasksService) deleteMemory(id string, subject string) bool {
	task, exists := s.memoryTasks[id]
	if !exists || task.Subject != subject {
		return false
	}
	delete(s.memoryTasks, id)
	return true
}

// SearchByTitle
func (s *TasksService) SearchByTitle(term string, subject string) ([]models.Task, error) {
	if s.useDatabase {
		return s.repo.SearchByTitle(term, subject)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []models.Task
	for _, task := range s.memoryTasks {
		if task.Subject == subject && contains(task.Title, term) {
			results = append(results, task)
		}
	}

	return results, nil
}

// Демонстрация SQL-инъекции
func (s *TasksService) SearchByTitleVulnerable(term string, subject string) ([]models.Task, error) {
	if s.useDatabase {
		s.log.Warn("Using VULNERABLE search method - FOR DEMO ONLY",
			zap.String("term", term),
			zap.String("subject", subject))
		return s.repo.SearchByTitleVulnerable(term, subject)
	}
	return s.SearchByTitle(term, subject)
}

// Улучшенная генерация ID для in-memory режима
func (s *TasksService) generateID() string {
	s.counter++
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("t%s-%05d", timestamp, s.counter)
}

func generateUUID() string {
	return uuid.New().String()
}

func contains(s, substr string) bool {
	if substr == "" {
		return true
	}
	if len(substr) > len(s) {
		return false
	}
	sLower := strings.ToLower(s)
	substrLower := strings.ToLower(substr)
	return strings.Contains(sLower, substrLower)
}
