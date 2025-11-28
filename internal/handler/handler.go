package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	handlerDto "github.com/berduk-dev/blog/internal/handler/dto"
	"github.com/berduk-dev/blog/internal/model"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreatePost(ctx context.Context, userID int, post handlerDto.CreatePostReq) error
	CreateComment(ctx context.Context, userID int, PostID int, comment handlerDto.CreateCommentReq) error
	GetPost(ctx context.Context, postID int) (model.Post, error)
	GetCommentsByPostID(ctx context.Context, postID int) ([]model.Comment, error)
	GetPostsByUserID(ctx context.Context, userID int) ([]model.Post, error)
	GetPosts(ctx context.Context) ([]model.Post, error)
}

type Handler struct {
	service Service
}

func New(service Service) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) CreatePost(c *gin.Context) {
	user := c.Value("user").(model.User)

	var req handlerDto.CreatePostReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		log.Println("error ShouldBindJSON:", err)
		return
	}

	err = h.service.CreatePost(c, user.ID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка! Попробуйте позже")
		log.Println("error service.CreatePost:", err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) CreateComment(c *gin.Context) {
	user := c.Value("user").(model.User)

	var req handlerDto.CreateCommentReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		log.Println("error ShouldBindJSON:", err)
		return
	}

	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		log.Println("error CreateComment, strconv.Atoi:", err)
		return
	}

	err = h.service.CreateComment(c, user.ID, postID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка! Попробуйте позже")
		log.Println("error service.CreateComment:", err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetPostsByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		log.Println("error GetPostsByUserID, strconv.Atoi:", err)
		return
	}

	posts, err := h.service.GetPostsByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка! Попробуйте позже")
		log.Println("error service.GetPostsByUserID:", err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetPosts(c *gin.Context) {
	posts, err := h.service.GetPosts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка! Попробуйте позже")
		log.Println("error service.GetPosts:", err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		log.Println("error GetPostsByUserID, strconv.Atoi:", err)
		return
	}

	post, err := h.service.GetPost(c, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка! Попробуйте позже")
		log.Println("error service.GetPosts:", err)
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) GetCommentsByPostID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		log.Println("error GetCommentsByPostID, strconv.Atoi:", err)
		return
	}

	comments, err := h.service.GetCommentsByPostID(c, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка! Попробуйте позже")
		log.Println("error service.GetCommentsByPostID:", err)
		return
	}

	c.JSON(http.StatusOK, comments)
}
