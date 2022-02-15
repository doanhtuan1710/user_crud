package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"user_crud/internal/pkg/infra"
	"user_crud/internal/pkg/repo"
	"user_crud/internal/pkg/setting"

	"github.com/tikivn/ultrago/u_logger"
)

type redisRepo struct {
	rd         *infra.RedisClient
	expiration time.Duration
}

func NewRedisRepo(rd *infra.RedisClient, expiration time.Duration) (r repo.UserRedisRepo) {

	r = &redisRepo{
		rd:         rd,
		expiration: expiration,
	}

	return
}

func (r *redisRepo) redisKey(key string) (res string) {
	res = fmt.Sprintf("%v:%v", setting.REDIS_APP_PREFIX, key)
	return
}

func (r *redisRepo) redisValue(value interface{}) (res string, err error) {

	logger := u_logger.NewLogger()

	data, err := json.Marshal(value)
	if err != nil {
		logger.Errorf("failed to redis value to json: %v", err)
		return
	}

	res = string(data)
	return
}

func (r *redisRepo) Set(ctx context.Context, key string, value interface{}) (err error) {
	redisKey := r.redisKey(key)
	redisValue, err := r.redisValue(value)
	if err != nil {
		return
	}
	err = r.rd.Set(ctx, redisKey, redisValue, r.expiration).Err()
	return
}

func (r *redisRepo) Get(ctx context.Context, key string) (res string, err error) {
	redisKey := r.redisKey(key)
	res, err = r.rd.Get(ctx, redisKey).Result()
	return
}
