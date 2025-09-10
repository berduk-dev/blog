package handler

import (
	"fmt"
	"github.com/berduk-dev/blog/internal/model"
	"github.com/berduk-dev/blog/internal/service"
	"github.com/gin-gonic/gin"
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
		return
	}

	err = h.service.CreateUser(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Попробуйте позже")
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("error GetUser strconv.Atoi: %w", err))
		return
	}

	user, err := h.service.GetUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("error service.GetUser: %w", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":    user.UserID,
		"name":       user.Name,
		"email":      user.Email,
		"is_admin":   user.IsAdmin,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}
