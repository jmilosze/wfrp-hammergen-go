package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"net/http"
)

func RegisterUserRoutes(router *gin.Engine, userService domain.UserService, jwtService domain.JwtService) {
	router.GET("api/user/:userId", RequireJwt(jwtService), getUserHandler(userService))
	router.POST("api/user", createUserHandler(userService))
}

func getUserHandler(userService domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		authUserId := c.GetString("authUserId")

		if userId != authUserId {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "unauthorized"})
			return
		}

		user, err := userService.FindUserById(userId)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": gin.H{"id": user.Id, "username": user.Username}})
	}
}

func createUserHandler(userService domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {

		var userData struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusNotFound, "message": err.Error()})
			return
		}

		user, err := userService.CreateUser(userData.Username, userData.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": gin.H{"id": user.Id, "username": user.Username}})
	}
}
