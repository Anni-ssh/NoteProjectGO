package redisDB

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisClient(cfg Config) (*redis.Client, error) {
	const op = "redisDB.NewRedisClient"
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redisDB client error: %w, operation: %s", err, op)
	}

	return client, nil
}
