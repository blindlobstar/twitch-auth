package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const PREFIX = "token::%s"

type TokenStore interface {
	SaveTokens(ctx context.Context, at string, rt string) error
	GetToken(ctx context.Context, rt string) (string, error)
}

type RedisTokenStore struct {
	RDB *redis.Client
}

func (ts *RedisTokenStore) SaveTokens(ctx context.Context, at string, rt string) error {

	key := fmt.Sprintf(PREFIX, rt)
	_, err := ts.RDB.SetNX(ctx, key, at, 5*time.Hour).Result()
	return err
}

func (ts *RedisTokenStore) GetToken(ctx context.Context, rt string) (string, error) {
	key := fmt.Sprintf(PREFIX, rt)
	return ts.RDB.Get(ctx, key).Result()
}
