package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/littleurl/api/internal/application"
)

func GatewayJsonResponse(code int, body interface{}) (*events.APIGatewayV2HTTPResponse, error) {
	jsonBody, err := json.Marshal(body)
	if body == nil {
		jsonBody = []byte("")
	}

	if err != nil {
		return nil, err
	}

	return &events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Body:       string(jsonBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		IsBase64Encoded: false,
	}, nil
}

func GatewayErrorResponse(code int, message string) *events.APIGatewayV2HTTPResponse {
	if message == "" {
		message = http.StatusText(code)
	}

	bodyBytes, _ := json.Marshal(application.ErrorResponseBody{
		StatusCode: code,
		Message:    message,
	})

	return &events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Body:       string(bodyBytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		IsBase64Encoded: false,
	}
}
