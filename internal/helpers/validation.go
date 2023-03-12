package helpers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"gitlab.com/deltabyte_/littleurl/api/internal/application"
)

func GatewayValidationResponse(validationError error) *events.APIGatewayV2HTTPResponse {
	// handle actual errors
	if _, ok := validationError.(*validator.InvalidValidationError); ok {
		return GatewayErrorResponse(500, "")
	}

	// marshal body
	validationErrors := validationError.Error()
	bodyBytes, _ := json.Marshal(application.ErrorResponseBody{
		StatusCode: 400,
		Message:    "ValidationFailed",
		Details:    validationErrors,
	})

	// create response
	return &events.APIGatewayV2HTTPResponse{
		StatusCode: 400,
		Body:       string(bodyBytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		IsBase64Encoded: false,
	}
}

func ValidateRequest(req interface{}) error {
	validate := validator.New()

	// run the validations
	if err := validate.Struct(req); err != nil {
		return err
	}

	return nil
}
