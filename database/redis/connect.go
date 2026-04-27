package redis

import (
	"aurora-agent/utils"
	"context"
	"fmt"
	"sync"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	client          *goredis.Client
	redisOnce       sync.Once
	redisConnectErr error
)

func RedisConnect() (*goredis.Client, error) {
	redisOnce.Do(func() {
		conf := utils.InitViper("conf", "redis", "yaml")
		host := conf.GetString("redis.host")
		port := conf.GetString("redis.port")
		username := conf.GetString("redis.username")
		password := conf.GetString("redis.password")
		db := conf.GetInt("redis.db")

		client = goredis.NewClient(&goredis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Username: username,
			Password: password,
			DB:       db,
		})

		if err := client.Ping(context.Background()).Err(); err != nil {
			redisConnectErr = err
			utils.Logger.Error("Failed to connect to redis.", zap.Error(err))
			return
		}

		utils.Logger.Info("Redis connected successfully.", zap.String("host", host), zap.String("port", port), zap.String("username", username))
	})

	return client, redisConnectErr
}

func Client() *goredis.Client {
	if client == nil {
		if _, err := RedisConnect(); err != nil {
			return nil
		}
	}
	return client
}
