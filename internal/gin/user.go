package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"net/http"
)

func RegisterUserRoutes(router *gin.Engine, userService domain.UserService) {
	router.GET("/user/:name", func(c *gin.Context) {
		user, _ := userService.FindUserById("0")
		c.String(http.StatusOK, "Hello %s", user.Username)
	})
}
