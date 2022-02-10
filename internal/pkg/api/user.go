package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"user_crud/internal/entity"
	"user_crud/internal/pkg/repo"
	"user_crud/pkg/util"

	"github.com/go-chi/chi/v5"
	"github.com/tikivn/ultrago/u_logger"
)

type UserHandler struct {
	repo repo.UserRepo
	*BaseHandler
}

func NewUserHandler(base *BaseHandler, repo repo.UserRepo) *UserHandler {
	return &UserHandler{repo: repo, BaseHandler: base}
}

func (h *UserHandler) Route() (mux chi.Router) {

	// Create multiplexing
	mux = chi.NewRouter()

	// CRUDL
	mux.Post("/create", h.Create)
	mux.Get("/retrieve/{id}", h.Retrieve)
	mux.Patch("/update/{id}", h.Update)
	mux.Delete("/delete/{id}", h.Delete)
	mux.Post("/list", h.List)

	return
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {

	var (
		ctx, logger = u_logger.GetLogger(r.Context())
		in          = new(entity.User)
	)

	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		logger.Error(err)
		h.badRequest(w, r, err)
		return
	}

	in.Password = util.HashPassword(in.Password)
	if err := h.repo.Create(ctx, in); err != nil {
		logger.Error(err)
		h.badRequest(w, r, err)
		return
	}

	h.success(w, r, in.ToResponse())
}

func (h *UserHandler) Retrieve(w http.ResponseWriter, r *http.Request) {

	var (
		ctx, logger = u_logger.GetLogger(r.Context())
		id          = chi.URLParam(r, "id")
	)

	if id == "" {
		err := errors.New("MissingParamErr")
		logger.Errorf("%v: missing id", err)
		h.badRequest(w, r, err)
		return
	}

	out, err := h.repo.Retrieve(ctx, id)
	if err != nil {
		logger.Errorf("%v: unable to get user", err)
		h.badRequest(w, r, err)
		return
	}

	h.success(w, r, out.ToResponse())
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {

	var (
		ctx, logger = u_logger.GetLogger(r.Context())
		id          = chi.URLParam(r, "id")
		in          = new(entity.User)
	)

	if id == "" {
		err := errors.New("MissingParamErr")
		logger.Errorf("%v: missing id", err)
		h.badRequest(w, r, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		logger.Error(err)
		h.badRequest(w, r, err)
		return
	}

	out, err := h.repo.Update(ctx, id, in)
	if err != nil {
		logger.Error(err)
		h.badRequest(w, r, err)
		return
	}

	h.success(w, r, out.ToResponse())
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {

	var (
		ctx, logger = u_logger.GetLogger(r.Context())
		id          = chi.URLParam(r, "id")
	)

	if id == "" {
		err := errors.New("MissingParamErr")
		logger.Errorf("%v: missing id", err)
		h.badRequest(w, r, err)
		return
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		logger.Errorf("%v: unable to delete user", err)
		h.badRequest(w, r, err)
		return
	}

	h.success(w, r, nil)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {

	var (
		ctx, logger = u_logger.GetLogger(r.Context())
		query       = new(entity.Query)
	)

	if err := json.NewDecoder(r.Body).Decode(query); err != nil && err != io.EOF {
		logger.Error(err)
		h.badRequest(w, r, err)
		return
	}

	users, err := h.repo.List(ctx, query)
	if err != nil {
		logger.Error(err)
		h.badRequest(w, r, err)
		return
	}

	out := make([]entity.UserResponse, len(users))
	for _, user := range users {
		out = append(out, *user.ToResponse())
	}

	h.success(w, r, out)
}
