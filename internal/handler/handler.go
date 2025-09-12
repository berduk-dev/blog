package handler

import (
	"github.com/berduk-dev/blog/internal/model"
	"github.com/berduk-dev/blog/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service service.Service
}

func New(service service.Service) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req model.CreateUserReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		log.Println("error ShouldBindJSON:", err)
		return
	}

	err = h.service.CreateUser(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка! Попробуйте позже")
		log.Println("error service.CreateUser:", err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) CreatePost(c *gin.Context) {
	var req model.CreatePostReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		log.Println("error ShouldBindJSON:", err)
		return
	}

	err = h.service.CreatePost(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка! Попробуйте позже")
		log.Println("error service.CreatePost:", err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) CreateComment(c *gin.Context) {
	var req model.CreateCommentReq

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

	err = h.service.CreateComment(c, postID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка! Попробуйте позже")
		log.Println("error service.CreateComment:", err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		log.Println("error GetUser, strconv.Atoi:", err)
		return
	}

	user, err := h.service.GetUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Произошла ошибка! Попробуйте позже")
		log.Println("error service.GetUser:", err)
		return
	}

	c.JSON(http.StatusOK, user)
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
