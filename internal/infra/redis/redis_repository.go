package repository

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewRedisClient(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return rdb
}

func BlockIP(rdb *redis.Client, ip string, blockTime time.Duration) error {
	key := "block:ip:" + ip
	return rdb.Set(context.Background(), key, 1, blockTime).Err()
}

func IsBlocked(rdb *redis.Client, ip string) (bool, error) {
	key := "block:ip:" + ip
	val, err := rdb.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return val == "1", nil
}
