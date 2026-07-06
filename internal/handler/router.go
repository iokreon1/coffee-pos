package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/iokreon1/coffee-pos/pkg/response"
)

func NewRouter(appEnv string) *gin.Engine {
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			response.OK(c, "server is running", nil)
		})
	}

	return r
}
