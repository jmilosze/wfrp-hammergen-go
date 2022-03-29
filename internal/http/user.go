package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterUserRoutes(router *gin.Engine) {
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})
}
