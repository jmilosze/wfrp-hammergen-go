package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"net/http"
)

func RegisterUserRoutes(router *gin.Engine, userService domain.UserService, jwtService domain.JwtService) {
	router.GET("api/user/:userId", RequireJwt(jwtService), getHandler(userService))
	router.DELETE("api/user/:userId", RequireJwt(jwtService), deleteHandler(userService))
	router.PUT("api/user/:userId", RequireJwt(jwtService), updateHandler(userService))
	router.POST("api/user", createHandler(userService))
}

func getHandler(userService domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		authUserId := c.GetString("authUserId")

		if userId != authUserId {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "unauthorized"})
			return
		}

		user, err := userService.GetById(userId)

		if err != nil {
			if err.Type == domain.UserNotFoundError {
				c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "user not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": gin.H{"id": user.Id, "username": user.Username}})
	}
}

func deleteHandler(userService domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		authUserId := c.GetString("authUserId")

		if userId != authUserId {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "unauthorized"})
			return
		}

		if err := userService.Delete(userId); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "internal server error"})
		}

		c.JSON(http.StatusNoContent, gin.H{"code": http.StatusNoContent})
	}
}

func createHandler(userService domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		var userData domain.User
		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusNotFound, "message": err.Error()})
			return
		}

		user, err := userService.Create(&userData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "data": gin.H{"id": user.Id, "username": user.Username}})
	}
}

func updateHandler(userService domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		authUserId := c.GetString("authUserId")

		if userId != authUserId {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "unauthorized"})
			return
		}

		var userData domain.User
		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusNotFound, "message": err.Error()})
			return
		}

		user, err := userService.Update(userId, &userData)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "internal server error"})
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": gin.H{"id": user.Id, "username": user.Username}})
	}
}
