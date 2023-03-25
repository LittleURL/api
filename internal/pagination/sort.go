package pagination

import "github.com/aws/aws-lambda-go/events"

const (
	SortDescQuery = "sortDesc"
	SortByQuery   = "sortBy"
)

func ParseSortDesc(event *events.APIGatewayV2HTTPRequest) bool {
	rawSortDesc := event.QueryStringParameters[SortDescQuery]

	sortDesc := false
	if rawSortDesc == "true" {
		sortDesc = true
	}

	return sortDesc
}

func ParseSortBy(event *events.APIGatewayV2HTTPRequest) string {
	return event.QueryStringParameters[SortByQuery]
}
