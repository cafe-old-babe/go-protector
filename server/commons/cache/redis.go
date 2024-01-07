package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
)

var _redis *redis.Client
var pingErr error
var once sync.Once

func InitRedis(options *redis.Options) error {
	once.Do(func() {
		_redis = redis.NewClient(options)
		pingErr = _redis.Ping(context.TODO()).Err()
	})

	return pingErr
}

func GetRedis() *redis.Client {
	return _redis
}
