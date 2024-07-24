package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
)

var _redis *redis.Client
var pingErr error
var once sync.Once

// 3-5	【实战】使用Redis存放图片验证码-掌握sync.Once函数与底层原理
func InitRedis(options *redis.Options) error {
	once.Do(func() {
		_redis = redis.NewClient(options)
		pingErr = _redis.Ping(context.TODO()).Err()
	})

	return pingErr
}

func GetRedisClient() *redis.Client {
	return _redis
}
