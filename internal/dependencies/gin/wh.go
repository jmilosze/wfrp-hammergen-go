package gin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"net/http"
)

func RegisterMutationRoutes(router *gin.Engine, ms domain.WHService[*domain.Mutation], js domain.JwtService) {
	router.POST("api/wh/mutation", RequireJwt(js), whCreateHandler[*domain.Mutation](ms))
}

func whCreateHandler[W domain.WhTypePointer](s domain.WHService[W]) func(*gin.Context) {
	return func(c *gin.Context) {
		var whData W
		if err := c.BindJSON(whData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		userId := c.Param("userId")
		whData.SetOwnerId(userId)

		whRead, _ := s.Create(c.Request.Context(), whData)

		returnData, err := whToMap(whRead)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "data": returnData})
	}
}

func whToMap[W domain.WhTypePointer](m W) (map[string]interface{}, error) {
	a, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	var res map[string]interface{}
	err = json.Unmarshal(a, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
