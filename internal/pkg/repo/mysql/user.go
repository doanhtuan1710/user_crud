package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"user_crud/internal/entity"
	"user_crud/internal/pkg/infra"
	"user_crud/internal/pkg/repo"
	"user_crud/internal/pkg/repo/redis"
	"user_crud/internal/pkg/setting"

	"gorm.io/gorm"
)

type userRepo struct {
	db  *gorm.DB
	rdb repo.UserRedisRepo
}

func NewUserRepo(db *gorm.DB, rd *infra.RedisClient) repo.UserRepo {
	return &userRepo{
		db:  db,
		rdb: redis.NewRedisRepo(rd, 1*time.Minute),
	}
}

func (r *userRepo) Create(ctx context.Context, in *entity.User) (err error) {

	err = r.db.WithContext(ctx).Create(in).Error
	return
}

func (r *userRepo) Retrieve(ctx context.Context, id string) (out *entity.User, err error) {

	out = new(entity.User)
	cache := new(userCache)

	// Check if data in redis
	res, err := r.rdb.Get(ctx, fmt.Sprintf(setting.REDIS_RETRIEVE_KEY, id))
	if err == nil {
		if errMarshall := json.Unmarshal([]byte(res), cache); errMarshall == nil {
			out = cache.User
			err = cache.GetError()
			return
		}
	}

	err = r.db.WithContext(ctx).First(out, id).Error

	// If not cache, then set cache
	defer func() {
		if err != nil {
			cache.Err = err.Error()
		}
		cache.User = out
		r.rdb.Set(ctx, fmt.Sprintf(setting.REDIS_RETRIEVE_KEY, id), cache)
	}()

	return
}

func (r *userRepo) Update(ctx context.Context, id string, in *entity.User) (out *entity.User, err error) {

	out = new(entity.User)
	if err = r.db.WithContext(ctx).Where("id = ?", id).First(out).Error; err != nil {
		err = fmt.Errorf("user with id %v not exists", id)
		return
	}

	if err = r.db.Model(out).Omit("id").Updates(in).Error; err != nil {
		err = fmt.Errorf("unable to update user data: %v", err)
		return
	}

	return
}

func (r *userRepo) Delete(ctx context.Context, id string) (err error) {

	err = r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
	return
}

func (r *userRepo) List(ctx context.Context, query *entity.Query) (out []*entity.User, err error) {

	out = make([]*entity.User, 0)

	qs := r.db.WithContext(ctx).
		Limit(query.GetLimit()).
		Offset(query.GetLimit() * (query.GetPage() - 1)).
		Order(fmt.Sprintf("%v %v", query.GetSortBy(), query.GetOrder()))

	for _, filter := range query.Filters {
		qs = qs.Where(fmt.Sprintf("%v %v ?", filter.GetKey(), filter.GetMethod()), filter.GetValue())
	}

	err = qs.Find(&out).Error

	return
}

type userCache struct {
	Err  string       `json:"err"`
	User *entity.User `json:"user"`
}

func (u *userCache) GetError() (err error) {

	if u.Err != "" {
		err = fmt.Errorf(u.Err)
	}

	return
}
