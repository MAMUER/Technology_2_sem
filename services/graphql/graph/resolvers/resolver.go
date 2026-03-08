package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

import (
	"tech-ip-sem2/services/graphql/internal/service"
	"tech-ip-sem2/shared/logger"
)

type Resolver struct {
	taskService *service.TaskService
	log         *logger.Logger
}

func NewResolver(taskService *service.TaskService, log *logger.Logger) *Resolver {
	return &Resolver{
		taskService: taskService,
		log:         log,
	}
}
