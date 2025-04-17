package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	Client *redis.Client
	TTL    time.Duration
}

func NewRedisClient(host, port string, ttl time.Duration) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// err := client.Ping(context.Background()).Err()
	err := client.Ping().Err()
	if err != nil {
		return nil, fmt.Errorf("try to ping to redis: %w", err)
	}
	c := &RedisClient{
		Client: client,
		TTL:    ttl,
	}
	return c, nil
}

func (rc *RedisClient) Close() error {
	return rc.Client.Close()
}
