package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"

	"github.com/nr-turkarslan/newrelic-tracing-golang/apps/proxy/commons"
	"github.com/nr-turkarslan/newrelic-tracing-golang/apps/proxy/services"
)

func CreateHandlers(
	router *gin.Engine,
	nrapp *newrelic.Application,
) {

	router.Use(nrgin.Middleware(nrapp))

	firstMethodService := services.FirstMethodService{}
	secondMethodService := services.SecondMethodService{
		Nrapp: nrapp,
	}
	thirdMethodService := services.ThirdMethodService{
		Nrapp:     nrapp,
		KafkaConn: createKafkaConnection(),
	}

	proxy := router.Group("/proxy")
	{
		// Health check
		proxy.GET("/health", func(ginctx *gin.Context) {
			ginctx.JSON(http.StatusOK, gin.H{
				"message": "OK!",
			})
		})

		// First method
		proxy.POST("/method1", firstMethodService.FirstMethod)

		// Second method
		proxy.POST("/method2", secondMethodService.SecondMethod)

		// Third method
		proxy.POST("/method3", thirdMethodService.ThirdMethod)
	}
}

func createKafkaConnection() *kafka.Conn {

	commons.Log(zerolog.InfoLevel, "Starting Kafka...")

	conn, err := kafka.DialLeader(context.Background(),
		"tcp", "kafka.kafka.svc.cluster.local:9092", "tracing", 0)
	if err != nil {
		commons.Log(zerolog.PanicLevel, err.Error())
		panic("could not dial: " + err.Error())
	}

	commons.Log(zerolog.InfoLevel, "Kafka is started.")
	return conn
}
