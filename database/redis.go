package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// Redis wraps a Redis client for use as a caching layer.
type Redis struct {
	client *redis.Client
}

// NewRedis creates a new Redis client. Addr must not be empty.
// Caller must call Close when done.
func NewRedis(addr string) *Redis {
	if addr == "" {
		return nil
	}
	client := redis.NewClient(&redis.Options{Addr: addr})
	return &Redis{client: client}
}

// NewRedisWithPing creates a Redis client and pings the server. Logs a warning on ping failure but returns the client.
func NewRedisWithPing(ctx context.Context, addr string) *Redis {
	r := NewRedis(addr)
	if r == nil {
		return nil
	}
	if err := r.Ping(ctx); err != nil {
		log.Printf("warning: redis ping failed: %v", err)
	}
	return r
}

// Ping checks connectivity to Redis.
func (r *Redis) Ping(ctx context.Context) error {
	if r == nil || r.client == nil {
		return nil
	}
	return r.client.Ping(ctx).Err()
}

// Close closes the Redis connection.
func (r *Redis) Close() error {
	if r == nil || r.client == nil {
		return nil
	}
	return r.client.Close()
}

// Client returns the underlying go-redis client for cache operations. May be nil.
func (r *Redis) Client() *redis.Client {
	if r == nil {
		return nil
	}
	return r.client
}
