package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"net/http"

	"golang.org/x/exp/slices"
)

func RegisterWhRoutes(router *gin.Engine) {
	router.GET("api/wh/:whType", whCreateHandler())
}

func whCreateHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		whType := c.Param("whType")
		if !slices.Contains(domain.WhTypes, whType) {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "invalid warhammer type"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "all goodie!"})
	}
}
