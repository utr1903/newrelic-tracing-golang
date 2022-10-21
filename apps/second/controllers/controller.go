package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	services "github.com/nr-turkarslan/newrelic-tracing-golang/apps/second/services"
)

func CreateHandlers(
	router *gin.Engine,
	nrapp *newrelic.Application,
) {

	router.Use(nrgin.Middleware(nrapp))
	proxy := router.Group("/second")
	{
		// Health check
		proxy.GET("/health", func(ginctx *gin.Context) {
			ginctx.JSON(http.StatusOK, gin.H{
				"message": "OK!",
			})
		})

		// Second method
		proxy.POST("/method2", services.SecondMethod)
	}
}
