package main

import (
	"context"
	"time"

	"log"

	"github.com/berduk-dev/blog/internal/cache"
	"github.com/berduk-dev/blog/internal/handler"
	"github.com/berduk-dev/blog/internal/middlewares"
	"github.com/berduk-dev/blog/internal/repo"
	"github.com/berduk-dev/blog/internal/service"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	connString := "postgres://postgres:password@localhost:5432/blog"
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД: ", err)
	}

	r := gin.Default()

	sessionsCache := cache.New(rdb)
	blogRepository := repo.New(conn)
	blogService := service.New(blogRepository, sessionsCache)
	blogHandler := handler.New(blogService)

	blogMiddlewares := middlewares.New(blogService)

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "User-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/register/", blogHandler.Register)
	r.POST("/login/", blogHandler.Login)
	// r.POST("/users/", blogHandler.CreateUser)
	r.GET("/users/:id/", blogHandler.GetUser)
	r.GET("/users/me", blogMiddlewares.SessionAuthMiddleware, blogHandler.GetCurrentUser)
	r.GET("/users/:id/posts/", blogHandler.GetPostsByUserID)
	r.POST("/posts/", blogMiddlewares.SessionAuthMiddleware, blogHandler.CreatePost)
	r.GET("/posts/", blogHandler.GetPosts)
	r.GET("/posts/:id/", blogHandler.GetPost)
	r.POST("/posts/:id/comments/", blogMiddlewares.SessionAuthMiddleware, blogHandler.CreateComment)
	r.GET("/posts/:id/comments/", blogHandler.GetCommentsByPostID)

	r.Run(":8088")
}
