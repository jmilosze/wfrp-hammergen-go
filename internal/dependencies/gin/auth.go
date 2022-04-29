package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"net/http"
	"strings"
)

func RegisterAuthRoutes(router *gin.Engine, userService domain.UserService, jwtService domain.JwtService) {
	router.POST("api/token", tokenHandler(userService, jwtService))
}

func tokenHandler(userService domain.UserService, jwtService domain.JwtService) func(*gin.Context) {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		user, err := userService.GetAndAuth(username, password)

		if err != nil {
			switch err.Type {
			case domain.UserNotFoundError:
				c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "user not found"})
			case domain.UserIncorrectPassword:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "invalid password"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			}
			return
		}

		var claims domain.Claims
		(&claims).Set(user)
		token, tokenErr := jwtService.GenerateToken(&claims)

		if tokenErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "error generating token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "access_token": token, "token_type": "bearer"})
	}
}

func RequireJwt(jwtService domain.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		token, err := parseAuthHeader(authHeader)
		if err != nil {
			unauthorized(c)
			return
		}

		claims, err := jwtService.ParseToken(token)
		if err != nil {
			unauthorized(c)
			return
		}

		c.Set("ClaimsId", claims.Id)
		c.Set("ClaimsAdmin", claims.Admin)
		c.Set("ClaimsSharedAccounts", claims.SharedAccounts)
	}
}

func unauthorized(c *gin.Context) {
	c.Abort()
	c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "unauthorized"})
}

func parseAuthHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("missing 'Authorization' header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", fmt.Errorf("invalid 'Authorization' header")
	}

	return parts[1], nil
}
