package service

import (
	"context"
	"errors"
	"gin_study/webook/internal/domain"
	"gin_study/webook/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("邮箱或密码错误")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}
func (svc UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	// 找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return u, ErrInvalidUserOrPassword
	}
	if err != nil {
		return u, err
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return u, ErrInvalidUserOrPassword
	}
	return u, err
}

func (svc UserService) UpdateNonSensitiveInfo(ctx *gin.Context, user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)
}
