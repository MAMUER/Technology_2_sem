package cache

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"tech-ip-sem2/services/tasks/internal/models"
	"tech-ip-sem2/shared/logger"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisCache struct {
	client  *redis.Client
	log     *logger.Logger
	enabled bool

	// Настройки TTL
	baseTTL   time.Duration
	jitterMax time.Duration

	// Префиксы ключей
	taskKeyPrefix string
	listKeyPrefix string
}

type CacheConfig struct {
	Addr      string
	Password  string
	DB        int
	BaseTTL   int // секунды
	JitterMax int // секунды
}

func NewRedisCache(cfg CacheConfig, log *logger.Logger) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		PoolTimeout:  4 * time.Second,
	})

	// Подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Warn("Redis connection failed, cache disabled", zap.Error(err))
		return &RedisCache{
			client:  client,
			log:     log,
			enabled: false,
		}
	}

	log.Info("Redis connected successfully",
		zap.String("addr", cfg.Addr),
		zap.Int("base_ttl", cfg.BaseTTL),
	)

	return &RedisCache{
		client:        client,
		log:           log,
		enabled:       true,
		baseTTL:       time.Duration(cfg.BaseTTL) * time.Second,
		jitterMax:     time.Duration(cfg.JitterMax) * time.Second,
		taskKeyPrefix: "tasks:task:",
		listKeyPrefix: "tasks:list:",
	}
}

// Закрытие соединения
func (c *RedisCache) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// Проверка доступности
func (c *RedisCache) IsEnabled() bool {
	return c.enabled
}

// Генерация ключа для задачи
func (c *RedisCache) taskKey(id string) string {
	return c.taskKeyPrefix + id
}

// Генерация ключа для списка задач пользователя
func (c *RedisCache) listKey(subject string) string {
	return c.listKeyPrefix + subject
}

// Вычисление TTL с jitter
func (c *RedisCache) ttlWithJitter() time.Duration {
	if c.jitterMax <= 0 {
		return c.baseTTL
	}

	jitter := time.Duration(rand.Int63n(int64(c.jitterMax)))
	return c.baseTTL + jitter
}

// Получение задачи из кэша
func (c *RedisCache) GetTask(ctx context.Context, id string) (*models.Task, error) {
	if !c.enabled {
		return nil, nil
	}

	key := c.taskKey(id)
	data, err := c.client.Get(ctx, key).Bytes()

	if err == redis.Nil {
		// Miss
		c.log.Debug("Cache miss", zap.String("key", key))
		return nil, nil
	}

	if err != nil {
		// Ошибка Redis
		c.log.Warn("Redis get error", zap.Error(err), zap.String("key", key))
		return nil, err
	}

	// Hit
	var task models.Task
	if err := json.Unmarshal(data, &task); err != nil {
		c.log.Warn("Failed to unmarshal cached task", zap.Error(err), zap.String("key", key))
		// Удаление битых данных
		c.client.Del(ctx, key)
		return nil, nil
	}

	c.log.Debug("Cache hit", zap.String("key", key), zap.String("task_id", task.ID))
	return &task, nil
}

// Сохранение задачи в кэш
func (c *RedisCache) SetTask(ctx context.Context, task *models.Task) error {
	if !c.enabled || task == nil {
		return nil
	}

	key := c.taskKey(task.ID)
	data, err := json.Marshal(task)
	if err != nil {
		c.log.Warn("Failed to marshal task for cache", zap.Error(err))
		return err
	}

	ttl := c.ttlWithJitter()
	err = c.client.Set(ctx, key, data, ttl).Err()

	if err != nil {
		c.log.Warn("Redis set error", zap.Error(err), zap.String("key", key))
		return err
	}

	c.log.Debug("Task cached", zap.String("key", key), zap.Duration("ttl", ttl))
	return nil
}

// Удаление задачи из кэша
func (c *RedisCache) DeleteTask(ctx context.Context, id string) error {
	if !c.enabled {
		return nil
	}

	key := c.taskKey(id)
	err := c.client.Del(ctx, key).Err()

	if err != nil {
		c.log.Warn("Redis delete error", zap.Error(err), zap.String("key", key))
		return err
	}

	c.log.Debug("Task deleted from cache", zap.String("key", key))
	return nil
}

// Получение списка задач из кэша
func (c *RedisCache) GetTaskList(ctx context.Context, subject string) ([]models.Task, error) {
	if !c.enabled {
		return nil, nil
	}

	key := c.listKey(subject)
	data, err := c.client.Get(ctx, key).Bytes()

	if err == redis.Nil {
		c.log.Debug("Cache miss for list", zap.String("subject", subject))
		return nil, nil
	}

	if err != nil {
		c.log.Warn("Redis get error for list", zap.Error(err))
		return nil, err
	}

	var tasks []models.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		c.log.Warn("Failed to unmarshal cached task list", zap.Error(err))
		c.client.Del(ctx, key)
		return nil, nil
	}

	c.log.Debug("Cache hit for list", zap.String("subject", subject), zap.Int("count", len(tasks)))
	return tasks, nil
}

// Сохранение списка задач в кэш
func (c *RedisCache) SetTaskList(ctx context.Context, subject string, tasks []models.Task) error {
	if !c.enabled {
		return nil
	}

	key := c.listKey(subject)
	data, err := json.Marshal(tasks)
	if err != nil {
		c.log.Warn("Failed to marshal task list for cache", zap.Error(err))
		return err
	}

	ttl := c.ttlWithJitter()
	err = c.client.Set(ctx, key, data, ttl).Err()

	if err != nil {
		c.log.Warn("Redis set error for list", zap.Error(err))
		return err
	}

	c.log.Debug("Task list cached", zap.String("subject", subject), zap.Duration("ttl", ttl))
	return nil
}

// Удаление списка задач
func (c *RedisCache) DeleteTaskList(ctx context.Context, subject string) error {
	if !c.enabled {
		return nil
	}

	key := c.listKey(subject)
	err := c.client.Del(ctx, key).Err()

	if err != nil {
		c.log.Warn("Redis delete error for list", zap.Error(err))
		return err
	}

	c.log.Debug("Task list deleted from cache", zap.String("subject", subject))
	return nil
}
