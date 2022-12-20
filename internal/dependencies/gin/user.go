package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"net/http"
)

func RegisterUserRoutes(router *gin.Engine, us domain.UserService, js domain.JwtService, cs domain.CaptchaService) {
	router.POST("api/user", userCreateHandler(us, cs))
	router.GET("api/user/:userId", RequireJwt(js), userGetHandler(us))
	router.GET("api/user", RequireJwt(js), userGetHandler(us))
	router.GET("api/user/exists/:userName", RequireJwt(js), userGetExistsHandler(us))
	router.GET("api/user/list", RequireJwt(js), userListHandler(us))
	router.PUT("api/user/:userId", RequireJwt(js), userUpdateHandler(us))
	router.PUT("api/user/credentials/:userId", RequireJwt(js), userUpdateCredentialsHandler(us))
	router.PUT("api/user/claims/:userId", RequireJwt(js), userUpdateClaimsHandler(us))
	router.DELETE("api/user/:userId", RequireJwt(js), userDeleteHandler(us))
	router.POST("api/user/send_reset_password", resetSendPasswordHandler(us, cs))
	router.POST("api/user/reset_password", resetPasswordHandler(us))
}

type UserCreate struct {
	Username       string   `json:"username"`
	Password       string   `json:"password"`
	SharedAccounts []string `json:"sharedAccounts"`
	Captcha        string   `json:"captcha"`
}

func userCreateHandler(us domain.UserService, cs domain.CaptchaService) func(*gin.Context) {
	return func(c *gin.Context) {
		var userData UserCreate
		if err := c.BindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		remoteAddr := c.Request.RemoteAddr
		if !cs.Verify(userData.Captcha, remoteAddr) {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "captcha verification error"})
			return
		}

		userWriteCredentials := domain.UserWriteCredentials{Username: userData.Username, Password: userData.Password}
		userWrite := domain.UserWrite{SharedAccounts: userData.SharedAccounts}

		userRead, err := us.Create(c.Request.Context(), &userWriteCredentials, &userWrite)
		if err != nil {
			switch err.Type {
			case domain.UserAlreadyExistsError:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "user already exists"})
			case domain.UserInvalidArgumentsError:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			}
			return
		}

		c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "data": userToMap(userRead)})
	}
}

func userToMap(user *domain.User) map[string]interface{} {
	return gin.H{
		"id":             user.Id,
		"username":       user.Username,
		"sharedAccounts": user.SharedAccounts,
		"admin":          user.Admin,
		"createdOn":      user.CreatedOn,
		"lastAuthOn":     user.LastAuthOn,
	}
}

func userGetHandler(us domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		claims := getUserClaims(c)

		if userId == "" {
			userId = claims.Id
		}

		user, err := us.Get(c.Request.Context(), claims, userId)

		if err != nil {
			switch err.Type {
			case domain.UserNotFoundError:
				c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "user not found"})
			case domain.UserUnauthorizedError:
				c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "unauthorized"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": userToMap(user)})
	}
}

func getUserClaims(c *gin.Context) *domain.Claims {
	var claims domain.Claims

	claims.Id = c.GetString("ClaimsId")
	claims.Admin = c.GetBool("ClaimsAdmin")

	sharedAccountsRaw, _ := c.Get("ClaimsSharedAccounts")
	claims.SharedAccounts, _ = sharedAccountsRaw.([]string)

	return &claims
}

func userGetExistsHandler(us domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userName")

		exists, err := us.Exists(c.Request.Context(), userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": gin.H{"exists": exists}})
	}
}

func userListHandler(us domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		claims := getUserClaims(c)
		allUsers, err := us.List(c.Request.Context(), claims)
		if err != nil {
			switch err.Type {
			case domain.UserUnauthorizedError:
				c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "setAnonymous"})
			default:
				c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "internal server error"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": allUsers})
	}
}

func userUpdateHandler(users domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		claims := getUserClaims(c)

		var userData domain.UserWrite
		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		userRead, err := users.Update(c.Request.Context(), claims, userId, &userData)
		if err != nil {
			switch err.Type {
			case domain.UserNotFoundError:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "user not found"})
			case domain.UserInvalidArgumentsError:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			case domain.UserUnauthorizedError:
				c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "setAnonymous"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": userToMap(userRead)})
	}
}

type UserCredentials struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	CurrentPassword string `json:"currentPassword"`
}

func userUpdateCredentialsHandler(us domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		claims := getUserClaims(c)

		var uc UserCredentials
		if err := c.ShouldBindJSON(&uc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		userWriteCredentials := domain.UserWriteCredentials{Username: uc.Username, Password: uc.Password}

		userRead, err := us.UpdateCredentials(c.Request.Context(), claims, userId, uc.CurrentPassword, &userWriteCredentials)
		if err != nil {
			switch err.Type {
			case domain.UserNotFoundError:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "user not found"})
			case domain.UserInvalidArgumentsError:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			case domain.UserIncorrectPasswordError:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "incorrect password"})
			case domain.UserUnauthorizedError:
				c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "setAnonymous"})
			default:
				c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "internal server error"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": userToMap(userRead)})
	}
}

func userUpdateClaimsHandler(us domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		claims := getUserClaims(c)

		var userData domain.UserWriteClaims
		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		userRead, err := us.UpdateClaims(c.Request.Context(), claims, userId, &userData)
		if err != nil {
			switch err.Type {
			case domain.UserNotFoundError:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "user not found"})
			case domain.UserInvalidArgumentsError:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			case domain.UserUnauthorizedError:
				c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "setAnonymous"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": userToMap(userRead)})
	}
}

func userDeleteHandler(us domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		claims := getUserClaims(c)

		if err := us.Delete(c.Request.Context(), claims, userId); err != nil {
			switch err.Type {
			case domain.UserUnauthorizedError:
				c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "setAnonymous"})
			default:
				c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "internal server error"})
			}
			return
		}

		c.JSON(http.StatusNoContent, gin.H{"code": http.StatusNoContent})
	}
}

type UserSendResetPassword struct {
	Username string `json:"username"`
	Captcha  string `json:"captcha"`
}

func resetSendPasswordHandler(us domain.UserService, cs domain.CaptchaService) func(*gin.Context) {
	return func(c *gin.Context) {
		var userData UserSendResetPassword
		if err := c.BindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		remoteAddr := c.Request.RemoteAddr
		if !cs.Verify(userData.Captcha, remoteAddr) {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "captcha verification error"})
			return
		}

		err := us.SendResetPassword(c.Request.Context(), userData.Username)

		if err != nil {
			switch err.Type {
			case domain.UserInvalidArgumentsError:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			case domain.UserNotFoundError:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "user not found"})
			case domain.UserSendEmailError:
				c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			}
			return
		}

		c.JSON(http.StatusNoContent, gin.H{"code": http.StatusNoContent})
	}
}

type UserResetPassword struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func resetPasswordHandler(us domain.UserService) func(*gin.Context) {
	return func(c *gin.Context) {
		var userData UserResetPassword
		if err := c.BindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		if err := us.ResetPassword(c.Request.Context(), userData.Token, userData.Password); err != nil {
			switch err.Type {
			case domain.UserInternalError:
				c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			}
			return
		}
		c.JSON(http.StatusNoContent, gin.H{"code": http.StatusNoContent})
	}

}
