package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"net/http"
)

func RegisterUserRoutes(router *gin.Engine, userService domain.UserService, jwtService domain.JwtService) {
	router.POST("api/user", createHandler(userService))
	router.GET("api/user/:userId", RequireJwt(jwtService), getHandler(userService))
	router.GET("api/user", RequireJwt(jwtService), listHandler(userService))
	router.PUT("api/user/:userId", RequireJwt(jwtService), updateHandler(userService))
	router.PUT("api/user/credentials/:userId", RequireJwt(jwtService), updateCredentialsHandler(userService))
	router.PUT("api/user/claims/:userId", RequireJwt(jwtService), updateClaims(userService))
	router.DELETE("api/user/:userId", RequireJwt(jwtService), deleteHandler(userService))
}

type UserCreate struct {
	Username       string   `json:"username" binding:"required"`
	Password       string   `json:"password" binding:"required"`
	SharedAccounts []string `json:"shared_accounts"`
}

func createHandler(userService domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		var userData UserCreate
		if err := c.BindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusNotFound, "message": err.Error()})
			return
		}

		userWriteCredentials := domain.UserWriteCredentials{Username: userData.Username, Password: userData.Password}
		userWrite := domain.UserWrite{SharedAccounts: userData.SharedAccounts}

		userRead, err := userService.Create(&userWriteCredentials, &userWrite)
		if err != nil {
			if err.Type == domain.UserAlreadyExistsError {
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "user already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "data": userToMap(userRead)})
	}
}

func userToMap(user *domain.User) map[string]interface{} {
	return gin.H{"id": user.Id, "username": user.Username, "shared_accounts": user.SharedAccounts, "admin": user.Admin}
}

func getHandler(userService domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		if !authorizeGet(c, userId) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "unauthorized"})
			return
		}

		user, err := userService.Get(userId)

		if err != nil {
			if err.Type == domain.UserNotFoundError {
				c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "user not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": userToMap(user)})
	}
}

func authorizeGet(c *gin.Context, userId string) bool {
	claims := getUserClaims(c)
	return userId == claims.Id || claims.Admin
}

func getUserClaims(c *gin.Context) *domain.Claims {
	var claims domain.Claims

	claims.Id = c.GetString("ClaimsId")
	claims.Admin = c.GetBool("ClaimsAdmin")

	sharedAccountsRaw, _ := c.Get("ClaimsSharedAccounts")
	claims.SharedAccounts, _ = sharedAccountsRaw.([]string)

	return &claims
}

func listHandler(userService domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		allUsers, err := userService.List()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "internal server error"})
			return
		}

		visibleUsers := authorizeList(c, allUsers)

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": visibleUsers})
	}
}

func authorizeList(c *gin.Context, userList []*domain.User) []*domain.User {
	claims := getUserClaims(c)

	var visibleUsers []*domain.User
	if claims.Admin {
		visibleUsers = userList
	} else {
		for _, u := range userList {
			if claims.Id == u.Id {
				visibleUsers = append(visibleUsers, u)
				break
			}
		}
	}

	return visibleUsers
}

type UserUpdate struct {
	SharedAccounts []string `json:"shared_accounts"`
}

func updateHandler(userService domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		if !authorizeModify(c, userId) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "unauthorized"})
			return
		}

		var userData UserUpdate
		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		userWrite := (domain.UserWrite)(userData)
		userRead, err := userService.Update(userId, &userWrite)
		if err != nil {
			if err.Type == domain.UserNotFoundError {
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "user not found"})
				return
			}
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": userToMap(userRead)})
	}
}

func authorizeModify(c *gin.Context, userId string) bool {
	claims := getUserClaims(c)
	return userId == claims.Id
}

type UserCredentials struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	CurrentPassword string `json:"current_password" binding:"required"`
}

func updateCredentialsHandler(userService domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		if !authorizeModify(c, userId) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "unauthorized"})
			return
		}

		var uc UserCredentials
		if err := c.ShouldBindJSON(&uc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		userWriteCredentials := domain.UserWriteCredentials{Username: uc.Username, Password: uc.Password}

		userRead, err := userService.UpdateCredentials(userId, uc.CurrentPassword, &userWriteCredentials)
		if err != nil {
			switch err.Type {
			case domain.UserNotFoundError:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "user not found"})
			case domain.UserIncorrectPassword:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "incorrect password"})
			default:
				c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "internal server error"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": userToMap(userRead)})
	}
}

type UserClaims struct {
	Admin *bool `json:"admin"`
}

func updateClaims(userService domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		if !authorizeModifyClaims(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "unauthorized"})
			return
		}

		var userData UserClaims
		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		userWriteClaims := (domain.UserWriteClaims)(userData)
		userRead, err := userService.UpdateClaims(userId, &userWriteClaims)
		if err != nil {
			if err.Type == domain.UserNotFoundError {
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "user not found"})
				return
			}
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": userToMap(userRead)})
	}
}

func authorizeModifyClaims(c *gin.Context) bool {
	claims := getUserClaims(c)
	return claims.Admin
}

func deleteHandler(userService domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		if !authorizeModify(c, userId) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "unauthorized"})
			return
		}

		if err := userService.Delete(userId); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "internal server error"})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{"code": http.StatusNoContent})
	}
}
