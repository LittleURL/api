package pagination

import (
	"github.com/aws/aws-lambda-go/events"
)

const (
	TokenHeader string = "x-pagination-token"
	TokenQuery  string = "pageToken"
)

type Stringable interface {
	MarshalString() (string, error)
	UnmarshalString(string) error
}

func ParsePaginationToken(event *events.APIGatewayV2HTTPRequest, keys Stringable) (bool, error) {
	var token string = ""

	// header
	if tokenHeader, tokenHeaderExists := event.Headers[TokenHeader]; tokenHeaderExists {
		token = tokenHeader
	}

	// query param
	if tokenQuery, tokenQueryExists := event.QueryStringParameters[TokenQuery]; tokenQueryExists {
		token = tokenQuery
	}

	if token == "" {
		return false, nil
	}

	// decode
	if err := keys.UnmarshalString(token); err != nil {
		return false, err
	}

	return true, nil
}

func SetPaginationToken(event *events.APIGatewayV2HTTPResponse, keys Stringable) error {
	if keys == nil {
		return nil
	}

	header, err := keys.MarshalString()
	if err != nil {
		return err
	}

	// set header
	event.Headers[TokenHeader] = header
	return nil
}
