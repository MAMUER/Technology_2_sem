package consumer

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"tech-ip-sem2/shared/logger"
)

type Consumer struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	queue    string
	prefetch int
	log      *logger.Logger
	instance string
}

type ConsumerConfig struct {
	URL      string
	Queue    string
	Prefetch int
}

func NewConsumer(config ConsumerConfig, log *logger.Logger) (*Consumer, error) {
	// Подключение к RabbitMQ
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Создание канала
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Установка prefetch
	err = ch.Qos(
		config.Prefetch, // prefetch count
		0,               // prefetch size
		false,           // global
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set prefetch: %w", err)
	}

	// Объявление очереди
	_, err = ch.QueueDeclare(
		config.Queue, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	instance := os.Getenv("WORKER_ID")
	if instance == "" {
		instance = "worker-unknown"
	}

	log.Info("RabbitMQ consumer connected",
		zap.String("queue", config.Queue),
		zap.Int("prefetch", config.Prefetch),
		zap.String("instance", instance),
	)

	return &Consumer{
		conn:     conn,
		channel:  ch,
		queue:    config.Queue,
		prefetch: config.Prefetch,
		log:      log,
		instance: instance,
	}, nil
}

func (c *Consumer) Start() error {
	// Подписка на очередь
	msgs, err := c.channel.Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	c.log.Info("Consumer started, waiting for messages...",
		zap.String("instance", c.instance),
	)

	// Бесконечный цикл обработки сообщений
	go func() {
		for msg := range msgs {
			c.processMessage(msg)
		}
	}()

	return nil
}

func (c *Consumer) processMessage(msg amqp.Delivery) {
	startTime := time.Now()

	var event map[string]interface{}
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		c.log.Error("Failed to unmarshal message",
			zap.Error(err),
			zap.String("body", string(msg.Body)),
		)
		// Nack without requeue
		msg.Nack(false, false)
		return
	}

	// Извлечение основных полей
	eventType, _ := event["event"].(string)
	taskID, _ := event["task_id"].(string)
	title, _ := event["title"].(string)
	subject, _ := event["subject"].(string)
	requestID, _ := event["request_id"].(string)
	timestamp, _ := event["ts"].(string)

	// Логирование полученных событий
	c.log.Info("Received message",
		zap.String("event", eventType),
		zap.String("task_id", taskID),
		zap.String("title", title),
		zap.String("subject", subject),
		zap.String("request_id", requestID),
		zap.String("timestamp", timestamp),
		zap.String("instance", c.instance),
		zap.Int("body_size", len(msg.Body)),
	)

	// Симуляция обработки
	switch eventType {
	case "task.created":
		c.log.Info("Task created event processed", zap.String("task_id", taskID))
	case "task.updated":
		c.log.Info("Task updated event processed", zap.String("task_id", taskID))
	case "task.deleted":
		c.log.Info("Task deleted event processed", zap.String("task_id", taskID))
	default:
		c.log.Warn("Unknown event type", zap.String("event", eventType))
	}

	// Искусственная задержка для демонстрации prefetch
	time.Sleep(500 * time.Millisecond)

	// Ack
	if err := msg.Ack(false); err != nil {
		c.log.Error("Failed to ack message",
			zap.Error(err),
			zap.String("task_id", taskID),
		)
	} else {
		c.log.Debug("Message acknowledged",
			zap.String("task_id", taskID),
			zap.Duration("processing_time", time.Since(startTime)),
		)
	}
}

func (c *Consumer) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
