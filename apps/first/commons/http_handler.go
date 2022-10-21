package commons

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	dto "github.com/nr-turkarslan/newrelic-tracing-golang/apps/first/dtos"
)

func CreateSuccessfulHttpResponse(
	ginctx *gin.Context,
	httpStatusCode int,
	responseDto *dto.ResponseDto,
) {
	ginctx.JSON(httpStatusCode, responseDto)
}

func CreateFailedHttpResponse(
	ginctx *gin.Context,
	httpStatusCode int,
	message string,
) {
	log.Error(message)

	responseDto := dto.ResponseDto{
		Message: message,
	}

	ginctx.JSON(httpStatusCode, responseDto)
}
