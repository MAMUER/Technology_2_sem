package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"
	"tech-ip-sem2/services/worker/internal/consumer"
	"tech-ip-sem2/shared/logger"
)

func main() {
	log := logger.New("worker")
	defer log.Sync()

	workerID := os.Getenv("WORKER_ID")
	if workerID == "" {
		workerID = "worker-1"
	}

	log.Info("Worker service starting",
		zap.String("worker_id", workerID),
		zap.String("version", "1.0.0"),
	)

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL environment variable not set")
	}

	prefetch := 1
	if prefetchStr := os.Getenv("RABBITMQ_PREFETCH"); prefetchStr != "" {
		if val, err := strconv.Atoi(prefetchStr); err == nil && val > 0 {
			prefetch = val
		}
	}

	// Создаем consumer для заданий
	jobConsumer, err := consumer.NewJobConsumer(consumer.JobConsumerConfig{
		URL:           rabbitURL,
		Queue:         "task_jobs",
		RetryExchange: "task_jobs_dlx",
		RetryQueue:    "task_jobs_retry",
		DLQ:           "task_jobs_dlq",
		Prefetch:      prefetch,
	}, log)
	if err != nil {
		log.Fatal("Failed to create job consumer", zap.Error(err))
	}
	defer jobConsumer.Close()

	jobConsumer.SetInstanceID(workerID)

	if err := jobConsumer.Start(); err != nil {
		log.Fatal("Failed to start job consumer", zap.Error(err))
	}

	log.Info("Worker fully initialized and waiting for jobs...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Worker shutting down...")
	time.Sleep(1 * time.Second)
}
