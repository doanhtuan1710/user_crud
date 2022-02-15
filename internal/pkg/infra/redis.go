package infra

import (
	"context"
	"time"
	"user_crud/internal/pkg/setting"

	goredis "github.com/go-redis/redis/v8"
	"github.com/tikivn/ultrago/u_logger"
)

type RedisClient struct {
	*goredis.Client
}

func NewRedisClient() *RedisClient {
	logger := u_logger.NewLogger()
	c := goredis.NewClient(&goredis.Options{
		Addr:         setting.REDIS_ADDRESS,
		Password:     "",
		DB:           setting.REDIS_DB,
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
	})
	_, err := c.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatalf("error creating redis client: %s", err.Error())
	}
	return &RedisClient{c}
}
