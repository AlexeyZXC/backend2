package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func basicWork(client *RedisClient) error {
	const (
		mKey         = "basic_key"
		notExistKey  = "basic_key" + "_not_exist"
		notExistKey2 = "sure_not_exist"
	)
	keys := []string{mKey, notExistKey, notExistKey2}
	// comment it if you want data from previous launch
	/**/
	err := client.Del(context.Background(), keys...).Err()
	if err != nil {
		return err
	}
	/**/
	item, err := client.GetRecord(mKey)
	if err != nil {
		return err
	}
	fmt.Printf("FIRST GetRecord for key %q `%s`\n", mKey, item)
	ttl := 5 * time.Second
	// добавляет запись, https://redis.io/commands/set
	err = client.Set(context.Background(), mKey, 1, ttl).Err()
	if err != nil {
		return err
	}
	// just try to uncomment
	// time.Sleep(ttl)
	item, err = client.GetRecord(mKey)
	if err != nil {
		return err
	}
	fmt.Printf("SECOND GetRecord for key %q `%s`\n", mKey, item)
	// https://redis.io/commands/incrby
	var count int64 = 2
	err = client.IncrBy(context.Background(), mKey, count).Err()
	if err != nil {
		return err
	}
	fmt.Printf("INCR for key %q on value %v\n", mKey, count)
	item, err = client.GetRecord(mKey)
	if err != nil {
		return err
	}
	fmt.Printf("THIRD GetRecord for key %q `%s`\n", mKey, item)
	// https://redis.io/commands/decrby
	err = client.Decr(context.Background(), mKey).Err()
	if err != nil {
		return err
	}
	fmt.Printf("DECR for key %q\n", mKey)
	item, err = client.GetRecord(mKey)
	if err != nil {
		return err
	}
	fmt.Printf("FOURS GetRecord for key %q `%s`\n", mKey, item)
	err = client.Incr(context.Background(), notExistKey).Err()
	if err != nil {
		return err
	}
	fmt.Printf("INCR for key %q\n", notExistKey)
	item, err = client.GetRecord(notExistKey)
	if err != nil {
		return err
	}
	fmt.Printf("THIRD GetRecord for key %q `%s`\n", notExistKey, item)
	// https://redis.io/commands/mget
	result, err := client.MGet(context.Background(), keys...).Result()
	if err != nil {
		return err
	}
	log.Printf("MGET for keys %v, result: %v", keys, result)
	return nil
}

func main() {
	const (
		host = "localhost"
		port = "6379"
	)
	client, err := NewRedisClient(host, port)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	if err := basicWork(client); err != nil {
		log.Fatal(err)
	}
}

func withStructWork(client *RedisClient) error {
	jsonKey := "json_key"
	// comment it if you want data from previous launch
	/**/
	err := client.Del(context.Background(), []string{jsonKey}...).Err()
	if err != nil {
		return err
	}
	/**/
	type exampleStruct struct {
		FieldOne string `json:"field_one"`
		FieldTwo string `json:"field_two"`
	}
	s := exampleStruct{
		FieldOne: "one",
		FieldTwo: "two",
	}
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	ttl := 5 * time.Second
	err = client.Set(context.Background(), jsonKey, data, ttl).Err()
	if err != nil {
		return err
	}
	item, err := client.GetRecord(jsonKey)
	if err != nil {
		return err
	}
	fmt.Printf("GetRecord for key %q `%s`\n", jsonKey, item)
	e := new(exampleStruct)
	if err := json.Unmarshal(item, e); err != nil {
		return err
	}
	fmt.Printf("example struct: %+v\n", e)
	return nil
}
