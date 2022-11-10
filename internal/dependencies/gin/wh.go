package gin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"net/http"
)

func RegisterMutationRoutes(router *gin.Engine, ms domain.WhService, js domain.JwtService) {
	router.POST("api/wh/mutation", RequireJwt(js), whCreateHandler(ms, domain.WhTypeMutation))
	router.GET("api/wh/mutation/:whId", RequireJwt(js), whGetHandler(ms))
}

func whCreateHandler(s domain.WhService, whType int) func(*gin.Context) {
	return func(c *gin.Context) {
		var whWrite domain.Warhammer

		switch whType {
		case domain.WhTypeMutation:
			whWrite = &domain.Mutation{}
		case domain.WhTypeSpell:
			whWrite = &domain.Spell{}
		}

		reqData, err1 := c.GetRawData()
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err1.Error()})
			return
		}

		if err := whWrite.PopulateFromJson(reqData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		claims := getUserClaims(c)
		whRead, err2 := s.Create(c.Request.Context(), whWrite, claims)
		if err2 != nil {
			switch err2.WhType {
			case domain.WhInvalidArgumentsError:
				c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err2.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			}
			return
		}

		returnData, err3 := structToMap(whRead)
		if err3 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusInternalServerError, "message": err3.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "data": returnData})
	}
}

func structToMap(m any) (map[string]any, error) {
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

func whGetHandler(s domain.WhService) func(*gin.Context) {
	return func(c *gin.Context) {
		whId := c.Param("whId")
		claims := getUserClaims(c)

		wh, err1 := s.Get(c.Request.Context(), whId, claims)

		if err1 != nil {
			switch err1.ErrType {
			case domain.WhNotFoundError:
				c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "wh not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			}
			return
		}

		returnData, err2 := structToMap(wh)
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": returnData})
	}
}
