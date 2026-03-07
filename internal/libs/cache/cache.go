package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func New(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func Key(format string, parts ...any) string {
	return fmt.Sprintf(format, parts...)
}

func (c *Cache) Get(ctx context.Context, key string, dest any) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) Flush(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

func (c *Cache) Remember(ctx context.Context, key string, ttl time.Duration, dest any, fn func() (any, error)) error {
	err := c.Get(ctx, key, dest)
	if err == nil {
		return nil
	}
	if err != redis.Nil {
		return err
	}
	val, err := fn()
	if err != nil {
		return err
	}
	if err := c.Set(ctx, key, val, ttl); err != nil {
		return err
	}
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}
