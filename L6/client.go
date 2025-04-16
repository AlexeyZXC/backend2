package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(host, port string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("try to ping to redis: %w", err)
	}
	c := &RedisClient{
		Client: client,
	}
	return c, nil
}

func (c *RedisClient) Close() error {
	return c.Client.Close()
}

func (c *RedisClient) GetRecord(mkey string) ([]byte, error) {
	data, err := c.Get(context.Background(), mkey).Bytes()
	if err == redis.Nil {
		// we got empty result, it's not an error
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return data, nil
}
