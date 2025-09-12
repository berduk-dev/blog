package main

import (
	"context"
	"github.com/berduk-dev/blog/internal/handler"
	"github.com/berduk-dev/blog/internal/repo"
	"github.com/berduk-dev/blog/internal/service"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {

	connString := "postgres://admin:admin@localhost:5433/blog"
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД: ", err)
	}

	r := gin.Default()

	blogRepository := repo.New(conn)
	blogService := service.New(blogRepository)
	blogHandler := handler.New(blogService)

	r.POST("/users/", blogHandler.CreateUser)
	r.GET("/users/:id/", blogHandler.GetUser)
	r.GET("/users/:id/posts/", blogHandler.GetPostsByUserID)
	r.POST("/posts/", blogHandler.CreatePost)
	r.GET("/posts/", blogHandler.GetPosts)
	r.GET("/posts/:id/", blogHandler.GetPost)
	r.POST("/posts/:id/comments/", blogHandler.CreateComment)
	r.GET("/posts/:id/comments/", blogHandler.GetCommentsByPostID)

	r.Run(":8088")
}
