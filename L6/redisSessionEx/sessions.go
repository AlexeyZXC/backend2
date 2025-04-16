package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type Session struct {
	Login     string
	Useragent string
}

type SessionID struct {
	ID string
}

func (c *RedisClient) Create(in Session) (*SessionID, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("marshal session ID: %w", err)
	}
	id := SessionID{
		ID: uuid.NewV4().String(),
	}
	mkey := newRedisKey(id.ID)
	err = c.Set(context.Background(), mkey, data, c.TTL).Err()
	if err != nil {
		return nil, fmt.Errorf("redis: set key %q: %w", mkey, err)
	}
	return &id, nil
}

func (c *RedisClient) Check(in SessionID) (*Session, error) {
	mkey := newRedisKey(in.ID)
	data, err := c.GetRecord(mkey)
	if err != nil {
		return nil, fmt.Errorf("redis: get record by key %q: %w", mkey, err)
	} else if data == nil {
		// add here custom err handling
		return nil, nil
	}
	sess := new(Session)
	err = json.Unmarshal(data, sess)
	if err != nil {
		return nil, fmt.Errorf("unmarshal to session info: %w", err)
	}
	return sess, nil
}

func (c *RedisClient) Delete(in SessionID) error {
	mkey := newRedisKey(in.ID)
	err := c.Del(context.Background(), mkey).Err()
	if err != nil {
		return fmt.Errorf("redis: trying to delete value by key %q: %w", mkey, err)
	}
	return nil
}

func newRedisKey(sessionID string) string {
	return fmt.Sprintf("sessions: %s", sessionID)
}
