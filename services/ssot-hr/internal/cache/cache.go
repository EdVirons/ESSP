package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Rdb *redis.Client
	TTL time.Duration
}

func New(addr, password string, db int, ttl time.Duration) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Cache{Rdb: rdb, TTL: ttl}
}

func (c *Cache) Ping(ctx context.Context) error {
	return c.Rdb.Ping(ctx).Err()
}
