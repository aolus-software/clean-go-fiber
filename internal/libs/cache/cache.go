package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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

func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}

func (c *Cache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
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

func (c *Cache) RefreshTokenID(ctx context.Context, id uuid.UUID) (string, error) {
	var tokenID string
	err := c.Get(ctx, Key(RefreshTokenKey, id.String()), &tokenID)
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return tokenID, err
}

func (c *Cache) RevokeSessions(ctx context.Context, ids ...uuid.UUID) error {
	keys := make([]string, 0, len(ids)*2)
	for _, id := range ids {
		keys = append(keys, Key(UserInformationKey, id.String()), Key(RefreshTokenKey, id.String()))
	}
	return c.Delete(ctx, keys...)
}

func (c *Cache) ForgetUsers(ctx context.Context, ids ...uuid.UUID) error {
	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, Key(UserInformationKey, id.String()))
	}
	return c.Delete(ctx, keys...)
}
