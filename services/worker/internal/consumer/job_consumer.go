package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"tech-ip-sem2/services/worker/internal/models"
	"tech-ip-sem2/services/worker/internal/processor"
	"tech-ip-sem2/services/worker/internal/storage"
	"tech-ip-sem2/shared/logger"
)

type JobConsumer struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	queue         string
	retryExchange string
	retryQueue    string
	dlq           string
	prefetch      int
	log           *logger.Logger
	instance      string
	processed     *storage.ProcessedMessages
	processor     *processor.TaskProcessor
}

type JobConsumerConfig struct {
	URL           string
	Queue         string
	RetryExchange string
	RetryQueue    string
	DLQ           string
	Prefetch      int
}

func NewJobConsumer(config JobConsumerConfig, log *logger.Logger) (*JobConsumer, error) {
	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

	// Попытка подключиться с ретраями
	for i := range 5 {
		conn, err = amqp.Dial(config.URL)
		if err == nil {
			break
		}
		log.Warn(fmt.Sprintf("Failed to connect to RabbitMQ (attempt %d/5), retrying...", i+1),
			zap.Error(err))
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ after 5 attempts: %w", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Попытка объявить очереди
	err = declareQueues(ch, config, log)
	if err != nil {
		log.Warn("Failed to declare queues, will rely on publisher", zap.Error(err))
	}

	instance := "worker-1" // будет переопределено из env

	log.Info("RabbitMQ job consumer connected",
		zap.String("main_queue", config.Queue),
		zap.String("retry_queue", config.RetryQueue),
		zap.String("dlq", config.DLQ),
		zap.Int("prefetch", config.Prefetch),
	)

	return &JobConsumer{
		conn:          conn,
		channel:       ch,
		queue:         config.Queue,
		retryExchange: config.RetryExchange,
		retryQueue:    config.RetryQueue,
		dlq:           config.DLQ,
		prefetch:      config.Prefetch,
		log:           log,
		instance:      instance,
		processed:     storage.NewProcessedMessages(24 * time.Hour),
		processor:     processor.NewTaskProcessor(log),
	}, nil
}

func declareQueues(ch *amqp.Channel, config JobConsumerConfig, log *logger.Logger) error {
	// Проверка существования очередей
	_, err := ch.QueueDeclarePassive(
		config.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue %s does not exist", config.Queue)
	}

	log.Info("Queues already exist",
		zap.String("main_queue", config.Queue),
		zap.String("retry_queue", config.RetryQueue),
		zap.String("dlq", config.DLQ))

	return nil
}

func (c *JobConsumer) SetInstanceID(id string) {
	c.instance = id
}

func (c *JobConsumer) Start() error {
	msgs, err := c.channel.Consume(
		c.queue,
		"",
		false, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	c.log.Info("Job consumer started, waiting for jobs...",
		zap.String("instance", c.instance),
	)

	go func() {
		for msg := range msgs {
			c.processJob(msg)
		}
	}()

	return nil
}

func (c *JobConsumer) processJob(msg amqp.Delivery) {
	startTime := time.Now()
	log := c.log.With(zap.String("instance", c.instance))

	var job models.ProcessTaskJob
	if err := json.Unmarshal(msg.Body, &job); err != nil {
		log.Error("Failed to unmarshal job",
			zap.Error(err),
			zap.String("body", string(msg.Body)),
		)
		// Отправка в DLQ
		c.sendToDLQ(msg, "invalid_format")
		msg.Ack(false)
		return
	}

	if c.processed.IsProcessed(job.MessageID) {
		log.Info("Job already processed (idempotency)",
			zap.String("message_id", job.MessageID),
			zap.String("task_id", job.TaskID),
		)
		msg.Ack(false)
		return
	}

	log.Info("Received job",
		zap.String("job_type", job.JobType),
		zap.String("task_id", job.TaskID),
		zap.Int("attempt", job.Attempt),
		zap.String("message_id", job.MessageID),
	)

	// Обработка задания
	err := c.processor.ProcessTask(job.TaskID)

	// Симуляция случайных ошибок для демонстрации
	if contains(job.TaskID, "fail") {
		err = fmt.Errorf("simulated processing error")
	}

	if err == nil {
		// Успех
		c.processed.MarkProcessed(job.MessageID)
		log.Info("Job processed successfully",
			zap.String("task_id", job.TaskID),
			zap.Duration("processing_time", time.Since(startTime)),
		)
		msg.Ack(false)
		return
	}

	// Ошибка обработки
	log.Error("Job processing failed",
		zap.Error(err),
		zap.String("task_id", job.TaskID),
		zap.Int("attempt", job.Attempt),
	)

	// Проверка количества попыток
	if job.Attempt >= models.MaxAttempts {
		// Превышено максимальное число попыток
		log.Warn("Max attempts exceeded, sending to DLQ",
			zap.String("task_id", job.TaskID),
			zap.Int("attempts", job.Attempt),
		)
		c.sendToDLQ(msg, "max_attempts_exceeded")
		msg.Ack(false)
		return
	}

	// Retry
	job.Attempt++

	// Обновление тела сообщения
	newBody, _ := json.Marshal(job)

	err = c.channel.PublishWithContext(
		context.Background(),
		c.retryExchange,
		c.retryQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         newBody,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			MessageId:    job.MessageID,
			Headers: amqp.Table{
				"x-retry-count": job.Attempt,
			},
		})

	if err != nil {
		log.Error("Failed to send to retry queue",
			zap.Error(err),
			zap.String("task_id", job.TaskID),
		)

		msg.Nack(false, false)
		return
	}

	log.Info("Job sent to retry queue",
		zap.String("task_id", job.TaskID),
		zap.Int("next_attempt", job.Attempt),
	)

	msg.Ack(false)
}

func (c *JobConsumer) sendToDLQ(msg amqp.Delivery, reason string) {
	err := c.channel.PublishWithContext(
		context.Background(),
		c.retryExchange,
		"dlq",
		false,
		false,
		amqp.Publishing{
			ContentType:  msg.ContentType,
			Body:         msg.Body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			MessageId:    msg.MessageId,
			Headers: amqp.Table{
				"x-failure-reason": reason,
				"x-original-queue": msg.RoutingKey,
			},
		})

	if err != nil {
		c.log.Error("Failed to send to DLQ",
			zap.Error(err),
			zap.String("message_id", msg.MessageId),
		)
	}
}

func (c *JobConsumer) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
