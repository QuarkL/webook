package dao

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

// User 与数据表一一对应
type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	NickName string `gorm:"type=varchar(128)"`
	Birthday int64
	AboutMe  string `gorm:"type:varchar(4096)"`

	// 创建时间 毫秒数
	Ctime int64
	// 更新时间 毫秒数
	Utime int64
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	// 存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}
func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	u := User{}
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) UpdateById(ctx *gin.Context, entity User) error {
	return dao.db.WithContext(ctx).Model(&entity).Where("id=?", entity.Id).First(&entity).Updates(map[string]any{
		"utime":     time.Now().UnixMilli(),
		"nick_name": entity.NickName,
		"about_me":  entity.AboutMe,
		"birthday":  entity.Birthday,
	}).Error
}
