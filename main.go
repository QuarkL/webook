package main

import (
	"gin_study/webook/config"
	"gin_study/webook/internal/repository"
	"gin_study/webook/internal/repository/dao"
	"gin_study/webook/internal/service"
	"gin_study/webook/internal/web"
	"gin_study/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()
	u := initUser(db)
	u.RegisterRoutes(server)
	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 如果不加这个前端读不到jwt token
		ExposeHeaders: []string{"x-jwt-token"},
		// 是否允许携带Cookie之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	}))

	store, err := redis.NewStore(16, "tcp", config.Config.Redis.Addr, "", []byte("gF6eW9fP9yW6mN0yN9zX3oJ0iI4jK2aG"), []byte("tP4qF6gO8oP1yV5nP1vI1nY0vX9hG6jA"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store))
	//server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/signup").IgnorePaths("/users/login").CheckLogin())
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().IgnorePaths("/users/signup").IgnorePaths("/users/login").CheckLogin())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}
func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic("failed to connect database")
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
