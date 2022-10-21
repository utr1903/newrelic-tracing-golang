package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/nr-turkarslan/newrelic-tracing-golang/apps/proxy/commons"
	dto "github.com/nr-turkarslan/newrelic-tracing-golang/apps/proxy/dtos"
)

type FirstMethodService struct{}

func (s *FirstMethodService) FirstMethod(
	ginctx *gin.Context,
) {

	commons.LogWithContext(ginctx, zerolog.InfoLevel, "First method is triggered...")

	requestBody, err := s.parseRequestBody(ginctx)

	if err != nil {
		return
	}

	responseDtoFromFirstService, err := s.makeRequestToFirstService(ginctx,
		requestBody)

	if err != nil {
		return
	}

	commons.LogWithContext(ginctx, zerolog.InfoLevel, "First method is executed.")

	commons.CreateSuccessfulHttpResponse(ginctx, http.StatusOK,
		s.createResponseDto(responseDtoFromFirstService))
}

func (*FirstMethodService) parseRequestBody(
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

func (*FirstMethodService) makeRequestToFirstService(
	ginctx *gin.Context,
	requestDto *dto.RequestDto,
) (
	*dto.ResponseDto,
	error,
) {

	url := "http://first.first.svc.cluster.local:8080/first/method1"

	requestDtoInBytes, _ := json.Marshal(requestDto)

	httpResponse, err := http.Post(url, "application/json",
		bytes.NewBuffer(requestDtoInBytes))

	if err != nil {
		commons.CreateFailedHttpResponse(ginctx, http.StatusBadRequest,
			"Call to FirstService has failed.")

		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		commons.CreateFailedHttpResponse(ginctx, http.StatusBadRequest,
			"Call to FirstService has failed.")

		return nil, errors.New("call to first service has failed")
	}

	defer httpResponse.Body.Close()

	responseDtoInBytes, err := ioutil.ReadAll(httpResponse.Body)

	if err != nil {
		commons.CreateFailedHttpResponse(ginctx, http.StatusBadRequest,
			"Response from first service could not be parsed.")

		return nil, err
	}

	var responseDto dto.ResponseDto
	json.Unmarshal(responseDtoInBytes, &responseDto)

	commons.LogWithContext(ginctx, zerolog.InfoLevel, "Value retrieved: "+requestDto.Value)
	commons.LogWithContext(ginctx, zerolog.InfoLevel, "Tag retrieved: "+requestDto.Tag)

	return &responseDto, nil
}

func (*FirstMethodService) createResponseDto(
	data *dto.ResponseDto,
) *dto.ResponseDto {
	return &dto.ResponseDto{
		Message: "Succeeded.",
		Value:   data.Value,
		Tag:     data.Tag,
	}
}
