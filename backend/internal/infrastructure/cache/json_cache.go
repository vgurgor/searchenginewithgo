package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// SetJSON marshals v as JSON and stores it to Redis with given TTL.
func SetJSON(ctx context.Context, rdb *redis.Client, key string, v any, ttl time.Duration) error {
	if rdb == nil || key == "" || ttl <= 0 {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, b, ttl).Err()
}

// GetJSON attempts to read key from Redis and unmarshal JSON into out.
// Returns (true, nil) on cache hit, (false, nil) on miss, and (false, err) on error.
func GetJSON(ctx context.Context, rdb *redis.Client, key string, out any) (bool, error) {
	if rdb == nil || key == "" {
		return false, nil
	}
	res, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	if len(res) == 0 {
		return false, nil
	}
	if err := json.Unmarshal(res, out); err != nil {
		return false, err
	}
	return true, nil
}


