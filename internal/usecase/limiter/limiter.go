package limiter

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/rodrigoachilles/rate-limiter/configs/logger"
	"time"
)

type Limiter struct {
	RedisDB    *redis.Client
	IPLimit    int
	TokenLimit int
	BlockTime  time.Duration
}

func NewLimiter(rdb *redis.Client, ipLimit, tokenLimit int, blockTime time.Duration) *Limiter {
	return &Limiter{
		RedisDB:    rdb,
		IPLimit:    ipLimit,
		TokenLimit: tokenLimit,
		BlockTime:  blockTime,
	}
}

func (l *Limiter) getLimitKey(prefix, identifier string) string {
	return prefix + ":" + identifier
}

func (l *Limiter) AllowRequest(ctx context.Context, identifier string, limit int) (bool, error) {
	key := l.getLimitKey("req", identifier)
	currentCount, err := l.RedisDB.Get(ctx, key).Int()

	logger.Info(fmt.Sprintf("IP %s request %d/%d", identifier, currentCount, limit))

	if errors.Is(err, redis.Nil) {
		err := l.RedisDB.Set(ctx, key, 1, time.Second).Err()
		if err != nil {
			return false, err
		}
		return true, nil
	} else if err != nil {
		return false, err
	}

	if currentCount >= limit {
		return false, nil
	}

	err = l.RedisDB.Incr(ctx, key).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}
