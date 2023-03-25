package pagination

import (
	"net/url"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
)

const (
	LimitHeader string = "x-pagination-limit"
	LimitQuery  string = "pageLimit"
	LimitMin    int32  = 1
	LimitMax    int32  = 100
)

func ParsePaginationLimit(event *events.APIGatewayV2HTTPRequest) (*int32, error) {
	var limit *int32 = nil

	// header
	if limitHeader, limitHeaderExists := event.Headers[LimitHeader]; limitHeaderExists {
		parsed, err := strconv.Atoi(limitHeader)
		if err != nil {
			return nil, err
		}

		limit = aws.Int32(int32(parsed))
	}

	// query param
	if rawLimitQuery, limitQueryExists := event.QueryStringParameters[LimitQuery]; limitQueryExists {
		limitQuery, err := url.QueryUnescape(rawLimitQuery)
		if err != nil {
			return nil, err
		}

		parsed, err := strconv.Atoi(limitQuery)
		if err != nil {
			return nil, err
		}

		limit = aws.Int32(int32(parsed))
	}

	// min/max
	if limit != nil && *limit < LimitMin {
		limit = aws.Int32(LimitMin)
	}
	if limit != nil && *limit > LimitMax {
		limit = aws.Int32(LimitMax)
	}

	return limit, nil
}