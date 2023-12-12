package redis

import (
	"context"
	"log"
	"strconv"

	"github.com/wit-id/blueprint-backend-go/toolkit/db"

	"github.com/go-redis/redis/v8"
)

func NewRedisDatabase(opt *db.RedisOption) (*redis.Client, error) {
	portString := strconv.Itoa(opt.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     opt.Host + ":" + portString,
		Password: opt.Password,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Printf("failed connected to redis: %s, error: %s", opt.Host, err.Error())
	} else {
		log.Println("successfully connected to redis", opt.Host)
	}

	return rdb, nil
}
