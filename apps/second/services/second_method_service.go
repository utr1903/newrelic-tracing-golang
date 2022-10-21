package services

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/nr-turkarslan/newrelic-tracing-golang/apps/second/commons"
	dto "github.com/nr-turkarslan/newrelic-tracing-golang/apps/second/dtos"
)

func SecondMethod(
	ginctx *gin.Context,
) {

	commons.LogWithContext(ginctx, zerolog.InfoLevel, "Second method is triggered...")

	for key, element := range ginctx.Request.Header {
		fmt.Println("Key:", key, "=>", "Element:", element)
		commons.Log(zerolog.InfoLevel, "headerKey :"+key+" headerValue :"+element[0])
	}

	requestBody, err := parseRequestBody(ginctx)

	if err != nil {
		return
	}

	commons.LogWithContext(ginctx, zerolog.InfoLevel, "Second method is executed.")

	commons.CreateSuccessfulHttpResponse(ginctx, http.StatusOK,
		createResponseDto(requestBody))
}

func parseRequestBody(
	ginctx *gin.Context,
) (
	*dto.RequestDto,
	error,
) {
	var requestDto dto.RequestDto

	err := ginctx.BindJSON(&requestDto)

	if err != nil {
		commons.CreateFailedHttpResponse(ginctx, http.StatusBadRequest,
			"Request body could not be parsed.")

		return nil, err
	}

	commons.LogWithContext(ginctx, zerolog.InfoLevel, "Value provided: "+requestDto.Value)
	commons.LogWithContext(ginctx, zerolog.InfoLevel, "Tag provided: "+requestDto.Tag)

	return &requestDto, nil
}

func createResponseDto(
	data *dto.RequestDto,
) *dto.ResponseDto {
	return &dto.ResponseDto{
		Message: "Succeeded.",
		Value:   data.Value,
		Tag:     data.Tag,
	}
}
