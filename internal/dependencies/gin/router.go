package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/vearne/gin-timeout"
	"time"
)

func NewRouter(requestTimeout time.Duration) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(timeout.Timeout(timeout.WithTimeout(requestTimeout)))
	return router
}
