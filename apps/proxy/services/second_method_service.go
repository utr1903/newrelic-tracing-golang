package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"

	"github.com/nr-turkarslan/newrelic-tracing-golang/apps/proxy/commons"
	dto "github.com/nr-turkarslan/newrelic-tracing-golang/apps/proxy/dtos"
)

type SecondMethodService struct {
	Nrapp *newrelic.Application
}

func (s *SecondMethodService) SecondMethod(
	ginctx *gin.Context,
) {

	commons.LogWithContext(ginctx, zerolog.InfoLevel, "Second method is triggered...")

	requestBody, err := s.parseRequestBody(ginctx)

	if err != nil {
		return
	}

	responseDtoFromSecondService, err := s.makeRequestToSecondService(ginctx,
		requestBody)

	if err != nil {
		return
	}

	commons.LogWithContext(ginctx, zerolog.InfoLevel, "Second method is executed.")

	commons.CreateSuccessfulHttpResponse(ginctx, http.StatusOK,
		s.createResponseDto(responseDtoFromSecondService))
}

func (*SecondMethodService) parseRequestBody(
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

func (s *SecondMethodService) makeRequestToSecondService(
	ginctx *gin.Context,
	requestDto *dto.RequestDto,
) (
	*dto.ResponseDto,
	error,
) {

	secondUrl := "http://second.second.svc.cluster.local:8080/second/method2"

	customAttributes := map[string]string{"mycustomattributekey": "mycustomattributevalue"}

	httpResponse, err := commons.PerformPostRequest(secondUrl, ginctx,
		requestDto, customAttributes)

	if err != nil {
		commons.CreateFailedHttpResponse(ginctx, http.StatusBadRequest,
			"Call to SecondService has failed.")

		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		commons.CreateFailedHttpResponse(ginctx, http.StatusBadRequest,
			"Call to SecondService has failed.")

		return nil, errors.New("call to second service has failed")
	}

	defer httpResponse.Body.Close()

	responseDtoInBytes, err := ioutil.ReadAll(httpResponse.Body)

	if err != nil {
		commons.CreateFailedHttpResponse(ginctx, http.StatusBadRequest,
			"Response from second service could not be parsed.")

		return nil, err
	}

	var responseDto dto.ResponseDto
	json.Unmarshal(responseDtoInBytes, &responseDto)

	commons.LogWithContext(ginctx, zerolog.InfoLevel, "Value retrieved: "+requestDto.Value)
	commons.LogWithContext(ginctx, zerolog.InfoLevel, "Tag retrieved: "+requestDto.Tag)

	return &responseDto, nil
}

func (*SecondMethodService) createResponseDto(
	data *dto.ResponseDto,
) *dto.ResponseDto {
	return &dto.ResponseDto{
		Message: "Succeeded.",
		Value:   data.Value,
		Tag:     data.Tag,
	}
}
