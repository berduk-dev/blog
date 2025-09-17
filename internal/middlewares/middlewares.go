package middlewares

import (
	"log"
	"net/http"

	"github.com/berduk-dev/blog/internal/service"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	service service.Service
}

func New(service service.Service) Middleware {
	return Middleware{
		service: service,
	}
}

func (m *Middleware) SessionAuthMiddleware(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	found, userID, err := m.service.GetUserIDBySessionID(c, sessionID)
	if err != nil {
		log.Println("error service.GetUserIDBySessionID:", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !found {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := m.service.GetUser(c, userID)
	if err != nil {
		log.Println("error service.GetUser:", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", user)
}
