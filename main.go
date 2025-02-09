package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	err := rdb.Set(ctx, "username", "John", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "username").Result()
	if err != nil {
		panic(err)
	}
	println("username", val)
}
