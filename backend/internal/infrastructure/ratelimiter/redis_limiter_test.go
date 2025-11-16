package ratelimiter

import (
	"context"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestRedisLimiter_WithinLimit(t *testing.T) {
	mr, _ := miniredis.Run()
	defer mr.Close()
	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	rl := NewRedisLimiter(rc, true)
	ctx := context.Background()
	ok, err := rl.CheckLimit(ctx, "p1", 3)
	if err != nil || !ok {
		t.Fatalf("expected ok, err=%v", err)
	}
	_ = rl.RecordRequest(ctx, "p1")
	_ = rl.RecordRequest(ctx, "p1")
	ok, _ = rl.CheckLimit(ctx, "p1", 3)
	if !ok {
		t.Fatalf("should be within limit")
	}
}

func TestRedisLimiter_ExceedLimit(t *testing.T) {
	mr, _ := miniredis.Run()
	defer mr.Close()
	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	rl := NewRedisLimiter(rc, true)
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		_ = rl.RecordRequest(ctx, "p2")
	}
	ok, _ := rl.CheckLimit(ctx, "p2", 3)
	if ok {
		t.Fatalf("expected over limit")
	}
}

func TestRedisLimiter_Disabled(t *testing.T) {
	rc := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"}) // unreachable
	rl := NewRedisLimiter(rc, false)
	ctx := context.Background()
	ok, err := rl.CheckLimit(ctx, "p3", 1)
	if err != nil || !ok {
		t.Fatalf("disabled limiter should pass, err=%v", err)
	}
}
