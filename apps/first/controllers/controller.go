package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	services "github.com/nr-turkarslan/newrelic-tracing-golang/apps/first/services"
)

func CreateHandlers(
	router *gin.Engine,
) {

	proxy := router.Group("/first")
	{
		// Health check
		proxy.GET("/health", func(ginctx *gin.Context) {
			ginctx.JSON(http.StatusOK, gin.H{
				"message": "OK!",
			})
		})

		// First method
		proxy.POST("/method1", services.FirstMethod)
	}
}
