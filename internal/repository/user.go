package repository

import (
	"context"
	"gin_study/webook/internal/domain"
	"gin_study/webook/internal/repository/dao"
	"github.com/gin-gonic/gin"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{dao: dao}
}
func (r *UserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Email:    u.Email,
		Birthday: u.Birthday.UnixMilli(),
		NickName: u.NickName,
		AboutMe:  u.AboutMe,
		Password: u.Password,
	}
}
func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	// 在这里操作缓存
}
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, err
}

func (r *UserRepository) UpdateNonZeroFields(ctx *gin.Context, user domain.User) error {
	return r.dao.UpdateById(ctx, r.toEntity(user))
}
