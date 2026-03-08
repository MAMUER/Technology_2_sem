package resolvers

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"tech-ip-sem2/services/graphql/graph/generated"
	"tech-ip-sem2/services/graphql/graph/model"
	"tech-ip-sem2/services/graphql/internal/middleware"
)

// CreateTask is the resolver for the createTask field.
func (r *mutationResolver) CreateTask(ctx context.Context, input model.CreateTaskInput) (*model.Task, error) {
	r.log.Info("GraphQL mutation: createTask",
		zap.String("title", input.Title),
	)

	subject := middleware.GetSubject(ctx)

	task, err := r.taskService.CreateTask(input, subject)
	if err != nil {
		r.log.Error("failed to create task", zap.Error(err))
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

// UpdateTask is the resolver for the updateTask field.
func (r *mutationResolver) UpdateTask(ctx context.Context, id string, input model.UpdateTaskInput) (*model.Task, error) {
	r.log.Info("GraphQL mutation: updateTask",
		zap.String("id", id),
	)

	subject := middleware.GetSubject(ctx)

	task, err := r.taskService.UpdateTask(id, input, subject)
	if err != nil {
		r.log.Error("failed to update task", zap.Error(err), zap.String("id", id))
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	if task == nil {
		return nil, nil
	}

	return task, nil
}

// DeleteTask is the resolver for the deleteTask field.
func (r *mutationResolver) DeleteTask(ctx context.Context, id string) (bool, error) {
	r.log.Info("GraphQL mutation: deleteTask",
		zap.String("id", id),
	)

	subject := middleware.GetSubject(ctx)

	deleted, err := r.taskService.DeleteTask(id, subject)
	if err != nil {
		r.log.Error("failed to delete task", zap.Error(err), zap.String("id", id))
		return false, fmt.Errorf("failed to delete task: %w", err)
	}

	return deleted, nil
}

// Tasks is the resolver for the tasks field.
func (r *queryResolver) Tasks(ctx context.Context) ([]*model.Task, error) {
	r.log.Info("GraphQL query: tasks")

	subject := middleware.GetSubject(ctx)

	tasks, err := r.taskService.GetAllTasks(subject)
	if err != nil {
		r.log.Error("failed to get tasks", zap.Error(err))
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

// Task is the resolver for the task field.
func (r *queryResolver) Task(ctx context.Context, id string) (*model.Task, error) {
	r.log.Info("GraphQL query: task",
		zap.String("id", id),
	)

	subject := middleware.GetSubject(ctx)

	task, err := r.taskService.GetTaskByID(id, subject)
	if err != nil {
		r.log.Error("failed to get task", zap.Error(err), zap.String("id", id))
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
