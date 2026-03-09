package processor

import (
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"
	"tech-ip-sem2/shared/logger"
)

type TaskProcessor struct {
	log *logger.Logger
}

func NewTaskProcessor(log *logger.Logger) *TaskProcessor {
	return &TaskProcessor{
		log: log,
	}
}

func (p *TaskProcessor) ProcessTask(taskID string) error {
	// Симуляция длительной обработки
	processingTime := time.Duration(2+rand.Intn(4)) * time.Second
	time.Sleep(processingTime)

	// 20% случайных ошибок для демонстрации
	if rand.Float32() < 0.2 {
		return fmt.Errorf("random processing error for task %s", taskID)
	}

	p.log.Debug("Task processed successfully",
		zap.String("task_id", taskID),
		zap.Duration("processing_time", processingTime),
	)

	return nil
}
