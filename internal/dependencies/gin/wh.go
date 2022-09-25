package gin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"net/http"
)

func RegisterMutationRoutes(router *gin.Engine, ms domain.WHService[domain.Mutation], js domain.JwtService) {
	router.POST("api/wh/mutation", RequireJwt(js), whCreateHandler(ms))
}

func whCreateHandler[W domain.WhType](s domain.WHService[W]) func(*gin.Context) {
	return func(c *gin.Context) {
		claims := getUserClaims(c)
		ownerId := claims.Id

		var whData W
		if err := c.BindJSON(&whData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		if err := domain.SetOwnerId(&whData, ownerId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
			return
		}

		whRead, err1 := s.Create(c.Request.Context(), &whData, claims)
		if err1 != nil {
			switch err1.WhType {
			case domain.WhInvalidArgumentsError:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err1.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			}
			return
		}

		returnData, err2 := whToMap(whRead)
		if err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusInternalServerError, "message": err2.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "data": returnData})
	}
}

func whToMap[W domain.WhType](m *W) (map[string]any, error) {
	a, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	var res map[string]any
	err = json.Unmarshal(a, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
