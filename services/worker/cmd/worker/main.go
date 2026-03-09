package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"tech-ip-sem2/services/worker/internal/consumer"
	"tech-ip-sem2/shared/logger"
)

func main() {
	// Логгер
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

	// Конфигурация RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "${RABBITMQ_URL}"
	}

	queueName := os.Getenv("RABBITMQ_QUEUE")
	if queueName == "" {
		queueName = "${RABBITMQ_QUEUE}"
	}

	prefetch := 1
	if prefetchStr := os.Getenv("RABBITMQ_PREFETCH"); prefetchStr != "" {
		if val, err := time.ParseDuration(prefetchStr); err == nil {
			prefetch = int(val)
		}
	}

	// Создание consumer
	cons, err := consumer.NewConsumer(consumer.ConsumerConfig{
		URL:      rabbitURL,
		Queue:    queueName,
		Prefetch: prefetch,
	}, log)
	if err != nil {
		log.Fatal("Failed to create consumer", zap.Error(err))
	}
	defer cons.Close()

	// Запуск consumer
	if err := cons.Start(); err != nil {
		log.Fatal("Failed to start consumer", zap.Error(err))
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Worker shutting down...")
	time.Sleep(1 * time.Second)
}
