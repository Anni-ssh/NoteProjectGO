package redisDB

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
)

type Redis struct {
	*redis.Client
}

func NewRedisStorage(r *redis.Client) *Redis {
	return &Redis{r}
}

func (r *Redis) CreateSession(ctx context.Context, userID, token string, expiration time.Duration) error {
	return r.Client.Set(ctx, userID, token, expiration).Err()
}

func (r *Redis) CheckSession(ctx context.Context, userID, token string) error {
	result, err := r.Client.Get(ctx, userID).Result()
	if err != nil {
		return fmt.Errorf("error getting session: %v", err)
	}

	if result != token {
		return ErrInvalidToken
	}
	return nil
}
