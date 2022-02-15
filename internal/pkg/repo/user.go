package repo

import (
	"context"
	"user_crud/internal/entity"
)

type UserRepo interface {
	Create(ctx context.Context, in *entity.User) (err error)
	Retrieve(ctx context.Context, id string) (out *entity.User, err error)
	Update(ctx context.Context, id string, in *entity.User) (out *entity.User, err error)
	Delete(ctx context.Context, id string) (err error)
	List(ctx context.Context, query *entity.Query) (out []*entity.User, err error)
}

type UserRedisRepo interface {
	Set(ctx context.Context, key string, value interface{}) (err error)
	Get(ctx context.Context, key string) (res string, err error)
}
