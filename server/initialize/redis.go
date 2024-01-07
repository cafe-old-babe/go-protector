package initialize

import (
	"github.com/redis/go-redis/v9"
	"go-protector/server/commons/cache"
	"go-protector/server/commons/config"
)

func initCache() error {
	redisConfig := config.GetConfig().Redis

	return cache.InitRedis(&redis.Options{
		Addr:     redisConfig.Addr,
		Username: redisConfig.Username,
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	})
}
