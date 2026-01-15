package db

import (
	"context"
	"server/config"
	"server/utils"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func initRedisDB(cfg *config.Config) (*redis.Client, error) {
	dsn := cfg.GetRedisDSN()
	utils.Logger.Infof("use Redis DSN:%s", dsn)

	redisConfig := &redis.Options{
		Addr:     dsn,
		Password: cfg.REDIS_PASSWORD,
		DB:       cfg.REDIS_DB,
	}
	client := redis.NewClient(redisConfig)

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return client, nil
}
