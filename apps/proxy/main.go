package main

import (
	"github.com/gin-gonic/gin"

	"github.com/nr-turkarslan/newrelic-tracing-golang/apps/proxy/commons"
	controller "github.com/nr-turkarslan/newrelic-tracing-golang/apps/proxy/controllers"
)

const PORT string = ":8080"

func main() {
	nrapp := commons.CreateNewRelicAgent()

	router := gin.Default()
	controller.CreateHandlers(router, nrapp)
	router.Run(PORT)
}
