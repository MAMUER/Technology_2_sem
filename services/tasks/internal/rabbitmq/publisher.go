package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"tech-ip-sem2/services/tasks/internal/models"
	"tech-ip-sem2/shared/logger"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	log     *logger.Logger
}

type PublisherConfig struct {
	URL   string
	Queue string
}

func NewPublisher(config PublisherConfig, log *logger.Logger) (*Publisher, error) {
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

	log.Info("RabbitMQ publisher connected",
		zap.String("queue", config.Queue),
	)

	return &Publisher{
		conn:    conn,
		channel: ch,
		queue:   config.Queue,
		log:     log,
	}, nil
}

func (p *Publisher) PublishEvent(ctx context.Context, eventType string, task models.Task, requestID string) error {
	// Создание события
	event := map[string]interface{}{
		"event":      eventType,
		"task_id":    task.ID,
		"title":      task.Title,
		"subject":    task.Subject,
		"ts":         time.Now().Format(time.RFC3339),
		"request_id": requestID,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Публикация сообщения
	err = p.channel.PublishWithContext(ctx,
		"",      // exchange
		p.queue, // routing key (имя очереди)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // persistent сообщение
			Timestamp:    time.Now(),
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	p.log.Debug("Event published",
		zap.String("event", eventType),
		zap.String("task_id", task.ID),
		zap.String("queue", p.queue),
	)

	return nil
}

func (p *Publisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
