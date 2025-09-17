package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/berduk-dev/blog/internal/model"
	"github.com/berduk-dev/blog/internal/model/errs"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	var req model.CreateUserReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		log.Println("error ShouldBindJSON:", err)
		return
	}

	err = h.service.CreateUser(c, req)
	if err != nil {
		if errors.Is(err, errs.ErrorEmailAlreadyExists) {
			c.JSON(http.StatusInternalServerError, "Эта почта уже используется")
			return
		}

		c.JSON(http.StatusInternalServerError, "Произошла ошибка! Попробуйте позже")
		log.Println("error service.CreateUser:", err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) Login(c *gin.Context) {
	var req model.LoginReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "У вас невалидный запрос")
		log.Println("error ShouldBindJSON:", err)
		return
	}

	user, err := h.service.AuthenticateUser(c, req.Email, req.Password)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sessionID, err := h.service.CreateSession(c, user.ID)
	if err != nil {
		log.Println("error service.CreateSession:", err)
		c.JSON(http.StatusInternalServerError, "Ошибка. Попробуйте позже")
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"session_id",                // name
		sessionID,                   // value
		int(24*time.Hour.Seconds()), // maxAge
		"/",                         // path
		"",                          // domain
		false,                       // secure (только HTTPS)
		false,                       // httpOnly (недоступен для JS)
	)

	c.JSON(200, model.LoginResp{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
	user := c.Value("user").(model.User)
	c.JSON(http.StatusOK, model.LoginResp{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
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
