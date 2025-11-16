package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	Client        *redis.Client
	Enabled       bool
	Window        time.Duration
}

func NewRedisLimiter(client *redis.Client, enabled bool) *RedisLimiter {
	return &RedisLimiter{
		Client:  client,
		Enabled: enabled,
		Window:  time.Minute,
	}
}

func (r *RedisLimiter) key(providerID string) string {
	return fmt.Sprintf("rl:%s:%d", providerID, time.Now().UTC().Unix()/int64(r.Window.Seconds()))
}

// CheckLimit ensures current count < limit
func (r *RedisLimiter) CheckLimit(ctx context.Context, providerID string, limit int) (bool, error) {
	if !r.Enabled {
		return true, nil
	}
	k := r.key(providerID)
	cntStr, err := r.Client.Get(ctx, k).Result()
	if err == redis.Nil {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	var cnt int64
	fmt.Sscanf(cntStr, "%d", &cnt)
	return cnt < int64(limit), nil
}

// RecordRequest increments count with TTL window
func (r *RedisLimiter) RecordRequest(ctx context.Context, providerID string) error {
	if !r.Enabled {
		return nil
	}
	k := r.key(providerID)
	pipe := r.Client.TxPipeline()
	pipe.Incr(ctx, k)
	pipe.Expire(ctx, k, r.Window+5*time.Second)
	_, err := pipe.Exec(ctx)
	return err
}


