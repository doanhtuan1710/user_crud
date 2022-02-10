package mysql

import (
	"context"
	"fmt"
	"user_crud/internal/entity"
	"user_crud/internal/pkg/repo"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) repo.UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, in *entity.User) (err error) {

	err = r.db.WithContext(ctx).Create(in).Error
	return
}

func (r *userRepo) Retrieve(ctx context.Context, id string) (out *entity.User, err error) {

	out = new(entity.User)
	err = r.db.WithContext(ctx).First(out, id).Error
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
