package application

import (
	"errors"
	"fmt"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorResponseBody struct {
	StatusCode int         `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details"`
}

type RequestError struct {
	StatusCode int
	Message string
	Err error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func (r *RequestError) GatewayResponse() *events.APIGatewayV2HTTPResponse {
	bodyBytes, _ := json.Marshal(ErrorResponseBody{
		StatusCode: r.StatusCode,
		Message:    r.Err.Error(),
	})

	return &events.APIGatewayV2HTTPResponse{
		StatusCode: r.StatusCode,
		Body:       string(bodyBytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		IsBase64Encoded: false,
	}
}

func (r *RequestError) GatewayResponseWithCode(code int) *events.APIGatewayV2HTTPResponse {
	r.StatusCode = code
	return r.GatewayResponse()
}

func NewRequestError(code int, message string, err error) *RequestError {
	if err == nil {
		err = errors.New(message)
	}

	return &RequestError{
		StatusCode: code,
		Message: message,
		Err: err,
	}
}