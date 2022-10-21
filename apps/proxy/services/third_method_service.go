package services

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"

	"github.com/nr-turkarslan/newrelic-tracing-golang/apps/proxy/commons"
	dto "github.com/nr-turkarslan/newrelic-tracing-golang/apps/proxy/dtos"
)

type ThirdMethodService struct {
	Nrapp     *newrelic.Application
	KafkaConn *kafka.Conn
}

func (s *ThirdMethodService) ThirdMethod(
	ginctx *gin.Context,
) {

	// Start transaction
	txn := newrelic.FromContext(ginctx)

	commons.LogWithContext(ginctx, zerolog.InfoLevel, "Third method is triggered...")

	requestBody, err := s.parseRequestBody(ginctx)

	if err != nil {
		txn.End()
		return
	}

	responseDtoFromThirdService, err := s.publishToKafka(ginctx,
		requestBody, txn)

	if err != nil {
		txn.End()
		return
	}

	commons.LogWithContext(ginctx, zerolog.InfoLevel, "Third method is executed.")

	commons.CreateSuccessfulHttpResponse(ginctx, http.StatusOK,
		s.createResponseDto(responseDtoFromThirdService))

	txn.End()
}

func (*ThirdMethodService) parseRequestBody(
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

func (s *ThirdMethodService) publishToKafka(
	ginctx *gin.Context,
	requestDto *dto.RequestDto,
	txn *newrelic.Transaction,
) (
	*dto.ResponseDto,
	error,
) {

	commons.LogWithContext(ginctx, zerolog.InfoLevel, "Publishing Kafka message...")
	requestDtoInBytes, _ := json.Marshal(requestDto)

	// Get distributed tracing headers
	dtHeaders := http.Header{}
	txn.InsertDistributedTraceHeaders(dtHeaders)

	// Put W3C headers into Kafka message
	headers := []kafka.Header{}
	headers = append(headers, kafka.Header{
		Key:   "traceparent",
		Value: []byte(dtHeaders.Get("traceparent")),
	})
	headers = append(headers, kafka.Header{
		Key:   "tracestate",
		Value: []byte(dtHeaders.Get("tracestate")),
	})

	_, err := s.KafkaConn.WriteMessages(kafka.Message{
		Headers: headers,
		Value:   requestDtoInBytes,
	})

	if err != nil {
		commons.CreateFailedHttpResponse(ginctx, http.StatusBadRequest,
			"Message could not be published.")

		return nil, err
	}

	commons.LogWithContext(ginctx, zerolog.InfoLevel, "Kafka message is published.")

	responseDto := dto.ResponseDto{
		Message: "Message is published.",
	}

	return &responseDto, nil
}

func (*ThirdMethodService) createResponseDto(
	data *dto.ResponseDto,
) *dto.ResponseDto {
	return &dto.ResponseDto{
		Message: "Succeeded.",
		Value:   data.Value,
		Tag:     data.Tag,
	}
}
