package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

var ctx = context.Background()

func RedisConnection() {
	redisURI := "rediss://default:AVNS_binh8fQhuH-feq_KF0b@redis-3d0497cc-dimas-1144.a.aivencloud.com:14420"

	addr, err := redis.ParseURL(redisURI)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(addr)

	err = rdb.Set(ctx, "key", "hello world", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("The value of key is:", val)
}
