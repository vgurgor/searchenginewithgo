package cache

import (
	"fmt"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisURL string) (*redis.Client, error) {
	u, err := url.Parse(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	addr := u.Host
	var password string
	if u.User != nil {
		if pw, ok := u.User.Password(); ok {
			password = pw
		}
	}
	db := 0

	opts := &redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		PoolTimeout:  4 * time.Second,
	}
	return redis.NewClient(opts), nil
}


