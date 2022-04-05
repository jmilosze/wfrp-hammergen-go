package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"net/http"
)

func RegisterUserRoutes(router *gin.Engine, userService domain.UserService, jwtService domain.JwtService) {
	router.GET("api/user/:id", RequireJwt(jwtService), func(c *gin.Context) {
		userId := c.Param("id")
		user, err := userService.FindUserById(userId)

		if err != nil {
			c.String(http.StatusNotFound, "User not found.")
		} else {
			c.String(http.StatusOK, "Hello %s", user.Username)
		}
	})
}
