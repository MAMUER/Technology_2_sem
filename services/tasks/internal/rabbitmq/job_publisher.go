package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"tech-ip-sem2/services/tasks/internal/jobs"
	"tech-ip-sem2/shared/logger"
)

type JobPublisher struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	queue         string
	dlq           string
	retryQueue    string
	retryExchange string
	log           *logger.Logger
}

type JobPublisherConfig struct {
	URL           string
	Queue         string
	DLQ           string
	RetryQueue    string
	RetryExchange string
	RetryTTL      int32 // milliseconds
}

func NewJobPublisher(config JobPublisherConfig, log *logger.Logger) (*JobPublisher, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Dead Letter Exchange
	err = ch.ExchangeDeclare(
		config.RetryExchange, // name
		"direct",             // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare DLX: %w", err)
	}

	// Основная очередь с DLX
	_, err = ch.QueueDeclare(
		config.Queue, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		amqp.Table{
			"x-dead-letter-exchange":    config.RetryExchange,
			"x-dead-letter-routing-key": config.RetryQueue,
		},
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare main queue: %w", err)
	}

	// Retry очередь с TTL
	_, err = ch.QueueDeclare(
		config.RetryQueue, // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		amqp.Table{
			"x-dead-letter-exchange":    "", // вернется в основную очередь
			"x-dead-letter-routing-key": config.Queue,
			"x-message-ttl":             config.RetryTTL,
		},
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare retry queue: %w", err)
	}

	// DLQ
	_, err = ch.QueueDeclare(
		config.DLQ, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare DLQ: %w", err)
	}

	// Связка retry exchange с очередями
	err = ch.QueueBind(
		config.RetryQueue,    // queue
		config.RetryQueue,    // routing key
		config.RetryExchange, // exchange
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind retry queue: %w", err)
	}

	err = ch.QueueBind(
		config.DLQ,           // queue
		"dlq",                // routing key
		config.RetryExchange, // exchange
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind DLQ: %w", err)
	}

	log.Info("RabbitMQ job queues initialized",
		zap.String("main_queue", config.Queue),
		zap.String("retry_queue", config.RetryQueue),
		zap.String("dlq", config.DLQ),
		zap.Int32("retry_ttl_ms", config.RetryTTL),
	)

	return &JobPublisher{
		conn:          conn,
		channel:       ch,
		queue:         config.Queue,
		dlq:           config.DLQ,
		retryQueue:    config.RetryQueue,
		retryExchange: config.RetryExchange,
		log:           log,
	}, nil
}

func (p *JobPublisher) PublishJob(ctx context.Context, taskID, messageID string) error {
	job := jobs.ProcessTaskJob{
		JobType:   jobs.JobTypeProcessTask,
		TaskID:    taskID,
		Attempt:   1,
		MessageID: messageID,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	body, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	err = p.channel.PublishWithContext(ctx,
		"",      // exchange
		p.queue, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			MessageId:    messageID,
		})
	if err != nil {
		return fmt.Errorf("failed to publish job: %w", err)
	}

	p.log.Debug("Job published",
		zap.String("task_id", taskID),
		zap.String("message_id", messageID),
		zap.String("queue", p.queue),
	)

	return nil
}

func (p *JobPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
