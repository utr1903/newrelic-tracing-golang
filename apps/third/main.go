package main

import (
	"github.com/gin-gonic/gin"

	controller "github.com/nr-turkarslan/newrelic-tracing-golang/apps/third/controllers"
)

const PORT string = ":8080"

func main() {
	router := gin.Default()
	controller.CreateHandlers(router)
	router.Run(PORT)
}
