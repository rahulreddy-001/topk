package cmd

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func getRedisClient(cfg *Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Database.Redis.Host, cfg.Database.Redis.Port),
		Password: cfg.Database.Redis.Password,
		DB:       cfg.Database.Redis.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
