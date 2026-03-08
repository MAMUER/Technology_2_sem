package service

import (
	// Убираем "time" - он не используется
	"github.com/google/uuid"
	"go.uber.org/zap"
	"tech-ip-sem2/services/graphql/graph/model"
	"tech-ip-sem2/services/graphql/internal/repository"
	"tech-ip-sem2/shared/logger"
)

type TaskService struct {
	repo repository.TaskRepository
	log  *logger.Logger
}

func NewTaskService(repo repository.TaskRepository, log *logger.Logger) *TaskService {
	return &TaskService{
		repo: repo,
		log:  log,
	}
}

func (s *TaskService) CreateTask(input model.CreateTaskInput, subject string) (*model.Task, error) {
	task := &model.Task{
		ID:          generateUUID(),
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
		Done:        false,
	}

	created, err := s.repo.Create(task, subject)
	if err != nil {
		s.log.Error("failed to create task", zap.Error(err))
		return nil, err
	}

	s.log.Info("task created via GraphQL",
		zap.String("task_id", created.ID),
		zap.String("title", created.Title),
	)

	return created, nil
}

func (s *TaskService) GetAllTasks(subject string) ([]*model.Task, error) {
	tasks, err := s.repo.GetAll(subject)
	if err != nil {
		s.log.Error("failed to get tasks", zap.Error(err))
		return nil, err
	}

	s.log.Debug("tasks listed via GraphQL", zap.Int("count", len(tasks)))
	return tasks, nil
}

func (s *TaskService) GetTaskByID(id string, subject string) (*model.Task, error) {
	task, err := s.repo.GetByID(id, subject)
	if err != nil {
		s.log.Error("failed to get task", zap.Error(err), zap.String("id", id))
		return nil, err
	}

	if task != nil {
		s.log.Debug("task retrieved via GraphQL", zap.String("task_id", id))
	}

	return task, nil
}

func (s *TaskService) UpdateTask(id string, input model.UpdateTaskInput, subject string) (*model.Task, error) {
	task, err := s.repo.Update(id, input, subject)
	if err != nil {
		s.log.Error("failed to update task", zap.Error(err), zap.String("id", id))
		return nil, err
	}

	if task != nil {
		s.log.Info("task updated via GraphQL", zap.String("task_id", id))
	}

	return task, nil
}

func (s *TaskService) DeleteTask(id string, subject string) (bool, error) {
	deleted, err := s.repo.Delete(id, subject)
	if err != nil {
		s.log.Error("failed to delete task", zap.Error(err), zap.String("id", id))
		return false, err
	}

	if deleted {
		s.log.Info("task deleted via GraphQL", zap.String("task_id", id))
	}

	return deleted, nil
}

func generateUUID() string {
	return uuid.New().String()
}
