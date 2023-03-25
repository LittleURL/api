package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorResponseBody struct {
	StatusCode int         `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details"`
}

type RequestError struct {
	StatusCode int
	Message    string
	Err        error
}

func (r *RequestError) ErrorMessage() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func (r *RequestError) Error() error {
	if r.StatusCode > 0 && r.Message != "" {
		fmt.Print(r.Err)
		return nil
	}

	return r.Err
}

func (r *RequestError) GatewayResponse() *events.APIGatewayV2HTTPResponse {
	if r.Message == "" {
		r.Message = http.StatusText(r.StatusCode)
	}

	bodyBytes, _ := json.Marshal(ErrorResponseBody{
		StatusCode: r.StatusCode,
		Message:    r.Message,
		Details:    r.Err,
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
		Message:    message,
		Err:        err,
	}
}
