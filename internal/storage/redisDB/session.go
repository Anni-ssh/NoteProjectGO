package redisDB

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
}

func NewRedisStorage(r *redis.Client) *Redis {
	return &Redis{r}
}

func (r *Redis) CreateSession(userID, token string) error {
	return r.Client.Set(context.Background(), userID, token, 0).Err()
}

func (r *Redis) CheckSession(userID, token string) error {
	result, err := r.Client.Get(context.Background(), userID).Result()
	if err != nil {
		return fmt.Errorf("error getting session: %v", err)
	}

	if result != token {
		return errors.New("wrong session token")
	}
	return nil
}
