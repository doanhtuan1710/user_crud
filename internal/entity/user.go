package entity

import (
	"strings"
)

type User struct {
	ID       int
	Name     string
	Password string
	Email    string
	Age      int
	Avatar   string
}

type UserResponse struct {
	ID     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Age    int    `json:"age,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

type FilterAttribute struct {
	Key    string `json:"key"`
	Method string `json:"method"`
	Value  string `json:"value"`
}

type Query struct {
	Limit   int                `json:"limit,omitempty"`
	Page    int                `json:"page,omitempty"`
	SortBy  string             `json:"sort_by,omitempty"`
	Order   string             `json:"order,omitempty"`
	Filters []*FilterAttribute `json:"filters,omitempty"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) ToResponse() (res *UserResponse) {

	res = &UserResponse{
		ID:     u.ID,
		Name:   u.Name,
		Email:  u.Email,
		Age:    u.Age,
		Avatar: u.Avatar,
	}

	return
}

func (q *Query) GetLimit() (limit int) {

	if q.Limit > 100 {
		limit = 100
	} else if q.Limit <= 0 {
		limit = 5
	} else {
		limit = q.Limit
	}

	return
}

func (q *Query) GetPage() (offset int) {

	if q.Page < 0 {
		offset = 0
	} else {
		offset = q.Page
	}

	return
}

func (q *Query) GetOrder() (order string) {

	order = strings.ToLower(strings.TrimSpace(q.Order))

	switch order {
	case "asc", "desc":
	default:
		order = "desc"
	}

	return
}

func (q *Query) GetSortBy() (sortBy string) {

	sortBy = strings.ToLower(strings.TrimSpace(q.SortBy))
	if sortBy == "" {
		sortBy = "id"
	}

	return
}

func (f *FilterAttribute) GetMethod() (method string) {

	method = strings.ToLower(strings.TrimSpace(f.Method))
	switch method {
	case ">", ">=", "<", "<=", "=":
	default:
		method = "="
	}

	return
}

func (f *FilterAttribute) GetKey() (key string) {

	key = strings.ToLower(strings.TrimSpace(f.Key))
	if key == "" {
		key = "id"
	}

	return
}

func (f *FilterAttribute) GetValue() (value string) {

	value = f.Value
	return
}
