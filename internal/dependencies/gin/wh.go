package gin

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain/warhammer"
)

func RegisterMutationRoutes(router *gin.Engine, ms warhammer.WhService, js domain.JwtService) {
	router.POST("api/wh/mutation", RequireJwt(js), whCreateOrUpdateHandler(true, ms, warhammer.WhTypeMutation))
	router.GET("api/wh/mutation/:whId", RequireJwt(js), whGetHandler(ms, warhammer.WhTypeMutation))
	router.PUT("api/wh/mutation/:whId", RequireJwt(js), whCreateOrUpdateHandler(false, ms, warhammer.WhTypeMutation))
	router.DELETE("api/wh/mutation/:whId", RequireJwt(js), whDeleteHandler(ms, warhammer.WhTypeMutation))
	router.GET("api/wh/mutation", RequireJwt(js), whListHandler(ms, warhammer.WhTypeMutation))
}

func RegisterSpellRoutes(router *gin.Engine, ms warhammer.WhService, js domain.JwtService) {
	router.POST("api/wh/spell", RequireJwt(js), whCreateOrUpdateHandler(true, ms, warhammer.WhTypeSpell))
	router.GET("api/wh/spell/:whId", RequireJwt(js), whGetHandler(ms, warhammer.WhTypeSpell))
	router.PUT("api/wh/spell/:whId", RequireJwt(js), whCreateOrUpdateHandler(false, ms, warhammer.WhTypeSpell))
	router.DELETE("api/wh/spell/:whId", RequireJwt(js), whDeleteHandler(ms, warhammer.WhTypeSpell))
	router.GET("api/wh/spell", RequireJwt(js), whListHandler(ms, warhammer.WhTypeSpell))
}

func RegisterPropertyRoutes(router *gin.Engine, ms warhammer.WhService, js domain.JwtService) {
	router.POST("api/wh/property", RequireJwt(js), whCreateOrUpdateHandler(true, ms, warhammer.WhTypeProperty))
	router.GET("api/wh/property/:whId", RequireJwt(js), whGetHandler(ms, warhammer.WhTypeProperty))
	router.PUT("api/wh/property/:whId", RequireJwt(js), whCreateOrUpdateHandler(false, ms, warhammer.WhTypeProperty))
	router.DELETE("api/wh/property/:whId", RequireJwt(js), whDeleteHandler(ms, warhammer.WhTypeProperty))
	router.GET("api/wh/property", RequireJwt(js), whListHandler(ms, warhammer.WhTypeProperty))
}

func whCreateOrUpdateHandler(isCreate bool, s warhammer.WhService, t warhammer.WhType) func(*gin.Context) {
	return func(c *gin.Context) {
		whWrite, err := warhammer.NewWh(t)
		if err != nil {
			c.JSON(ServerErrResp(""))
			return
		}

		reqData, err := c.GetRawData()
		if err != nil {
			c.JSON(BadRequestErrResp(err.Error()))
			return
		}

		if err = json.Unmarshal(reqData, &whWrite.Object); err != nil {
			c.JSON(BadRequestErrResp(err.Error()))
			return
		}

		claims := getUserClaims(c)

		var whRead *warhammer.Wh
		var whErr *warhammer.WhError
		if isCreate {
			whRead, whErr = s.Create(c.Request.Context(), t, &whWrite, claims)
		} else {
			whWrite.Id = c.Param("whId")
			whRead, whErr = s.Update(c.Request.Context(), t, &whWrite, claims)
		}

		if whErr != nil {
			switch whErr.ErrType {
			case warhammer.WhInvalidArgumentsError:
				c.JSON(BadRequestErrResp(whErr.Error()))
			case warhammer.WhUnauthorizedError:
				c.JSON(UnauthorizedErrResp(""))
			case warhammer.WhNotFoundError:
				c.JSON(NotFoundErrResp(""))
			default:
				c.JSON(ServerErrResp(""))
			}
			return
		}

		returnData, err := whToMap(whRead)
		if err != nil {
			c.JSON(ServerErrResp(""))
			return
		}

		c.JSON(OkResp(returnData))
	}
}

func whToMap(w *warhammer.Wh) (map[string]any, error) {
	whMap, err := structToMap(w.Object)
	if err != nil {
		return map[string]any{}, fmt.Errorf("error while mapping wh structure %s", err)
	}
	whMap["id"] = w.Id
	whMap["ownerId"] = w.OwnerId
	whMap["canEdit"] = w.CanEdit

	return whMap, nil
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

	if res == nil {
		res = map[string]any{}
	}

	return res, nil
}

func whGetHandler(s warhammer.WhService, t warhammer.WhType) func(*gin.Context) {
	return func(c *gin.Context) {
		whId := c.Param("whId")
		claims := getUserClaims(c)

		wh, whErr := s.Get(c.Request.Context(), t, whId, claims)

		if whErr != nil {
			switch whErr.ErrType {
			case warhammer.WhNotFoundError:
				c.JSON(NotFoundErrResp(""))
			default:
				c.JSON(ServerErrResp(""))
			}
			return
		}

		returnData, err := whToMap(wh)
		if err != nil {
			c.JSON(ServerErrResp(""))
			return
		}

		c.JSON(OkResp(returnData))
	}
}

func whListToListMap(whs []*warhammer.Wh) ([]map[string]any, error) {
	list := make([]map[string]any, len(whs))

	var err error
	for i, v := range whs {
		list[i], err = whToMap(v)
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}

func whDeleteHandler(s warhammer.WhService, t warhammer.WhType) func(*gin.Context) {
	return func(c *gin.Context) {
		whId := c.Param("whId")
		claims := getUserClaims(c)

		whErr := s.Delete(c.Request.Context(), t, whId, claims)

		if whErr != nil {
			switch whErr.ErrType {
			case warhammer.WhUnauthorizedError:
				c.JSON(UnauthorizedErrResp(""))
			default:
				c.JSON(ServerErrResp(""))
			}
			return
		}

		c.JSON(OkResp(""))
	}
}

func whListHandler(s warhammer.WhService, t warhammer.WhType) func(*gin.Context) {
	return func(c *gin.Context) {
		claims := getUserClaims(c)

		whs, whErr := s.List(c.Request.Context(), t, claims)

		if whErr != nil {
			c.JSON(ServerErrResp(""))
			return
		}

		returnData, err := whListToListMap(whs)
		if err != nil {
			c.JSON(ServerErrResp(""))
			return
		}

		c.JSON(OkResp(returnData))
	}
}
