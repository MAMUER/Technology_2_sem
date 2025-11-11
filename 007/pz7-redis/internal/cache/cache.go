package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	ctx context.Context
	rdb *redis.Client
}

func New(addr string) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	
	// Проверяем подключение к Redis
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
	} else {
		log.Printf("Successfully connected to Redis: %s", addr)
	}

	return &Cache{
		ctx: ctx,
		rdb: rdb,
	}
}

func (c *Cache) Set(key string, value string, ttl time.Duration) error {
	log.Printf("Setting key: %s=%s with TTL: %v", key, value, ttl)
	err := c.rdb.Set(c.ctx, key, value, ttl).Err()
	if err != nil {
		log.Printf("Error setting key %s: %v", key, err)
		return err
	}
	log.Printf("Successfully set key: %s", key)
	return nil
}

func (c *Cache) Get(key string) (string, error) {
	log.Printf("Getting key: %s", key)
	val, err := c.rdb.Get(c.ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Key not found: %s", key)
		return "", fmt.Errorf("nil")
	} else if err != nil {
		log.Printf("Error getting key %s: %v", key, err)
		return "", err
	}
	log.Printf("Retrieved key: %s=%s", key, val)
	return val, nil
}

func (c *Cache) TTL(key string) (time.Duration, error) {
	log.Printf("Getting TTL for key: %s", key)
	return c.rdb.TTL(c.ctx, key).Result()
}