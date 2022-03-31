package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine) {
	router.POST("api/token", tokenHandler)
}

func tokenHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	fmt.Printf("Username %s Password %s", username, password)
}
