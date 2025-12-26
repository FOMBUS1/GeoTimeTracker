package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(c *redis.Client) *Cache {
	return &Cache{client: c}
}

func (c *Cache) Get(ctx context.Context, lat, long float32) (string, error) {
	key := fmt.Sprintf("geo:%.5f:%.5f", lat, long)
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Set(ctx context.Context, lat, long float32, addr string) error {
	key := fmt.Sprintf("geo:%.5f:%.5f", lat, long)
	return c.client.Set(ctx, key, addr, 24*time.Hour).Err()
}
