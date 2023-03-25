package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"gitlab.com/deltabyte_/littleurl/api/internal/application"
	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
	"gitlab.com/deltabyte_/littleurl/api/internal/helpers"
	"gitlab.com/deltabyte_/littleurl/api/internal/pagination"
)

func Handler(ctx context.Context, event *events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	// init app
	app, err := application.New(ctx)
	if err != nil {
		fmt.Print("Failed to initialize app")
		panic(err)
	}

	// extract the user's ID
	userId, exists := event.RequestContext.Authorizer.JWT.Claims["sub"]
	if !exists {
		return helpers.GatewayErrorResponse(401, "UserID not found in auth token"), nil
	}

	// validate path params
	domainId, exists := event.PathParameters["domainId"]
	if !exists {
		return helpers.GatewayErrorResponse(400, "Missing `domainId` path parameter"), nil
	}

	// check permissions
	userRole, err := helpers.FindUserRole(app, domainId, userId)
	if err != nil || userRole == nil || !userRole.Role().LinksRead() {
		return helpers.GatewayErrorResponse(403, ""), err
	}

	// base query
	baseQueryExpression := []string{"domain_id = :domain"}
	baseQueryValues := map[string]ddbTypes.AttributeValue{
		":domain": &ddbTypes.AttributeValueMemberS{
			Value: domainId,
		},
	}

	// prefix query
	rawPrefix, prefixExists := event.QueryStringParameters["prefix"]
	if prefixExists {
		// decode
		prefix, err := url.QueryUnescape(rawPrefix)
		if err != nil {
			return helpers.GatewayErrorResponse(500, "Failed to parse prefix param"), err
		}

		// update query
		baseQueryExpression = append(baseQueryExpression, "begins_with(uri, :prefix)")
		baseQueryValues[":prefix"] = &ddbTypes.AttributeValueMemberS{
			Value: prefix,
		}
	}

	// sorting query(s)
	sortDesc := pagination.ParseSortDesc(event)
	var queryIndex *string = nil
	switch pagination.ParseSortBy(event) {
	case "updated":
	case "updated_at":
		queryIndex = aws.String("updated")
	}

	// pagination limit
	paginationLimit, err := pagination.ParsePaginationLimit(event)
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to parse pagination limit"), err
	}

	// pagination token
	var paginationKeys map[string]ddbTypes.AttributeValue
	paginationKeyEntity := &entities.LinkKey{}
	paginationTokenExists, err := pagination.ParsePaginationToken(event, paginationKeyEntity)
	if err == nil {
		paginationKeys, err = paginationKeyEntity.MarshalDynamoAV()
	}
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to parse pagination token"), err
	}

	// get links
	queryInput := &dynamodb.QueryInput{
		TableName:                 &app.Cfg.Tables.Links,
		IndexName:                 queryIndex,
		Limit:                     paginationLimit,
		ScanIndexForward:          aws.Bool(!sortDesc),
		KeyConditionExpression:    aws.String(strings.Join(baseQueryExpression, " AND ")),
		ExpressionAttributeValues: baseQueryValues,
	}
	if paginationTokenExists {
		queryInput.ExclusiveStartKey = paginationKeys
	}
	linksRes, err := app.DDBClient.Query(app.Ctx, queryInput)
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to load links"), err
	}

	// unmarshal links
	links := entities.Links{}
	if err := links.UnmarshalDynamoAV(linksRes.Items); err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to unmarshal links"), err
	}

	// response
	res, err := helpers.GatewayJsonResponse(200, links)

	// pagination
	if len(linksRes.LastEvaluatedKey) > 0 {
		linkKey := &entities.LinkKey{}

		// encode ddb to string
		if err := linkKey.UnmarshalDynamoAV(linksRes.LastEvaluatedKey); err != nil {
			return helpers.GatewayErrorResponse(500, "Failed to encode pagination token"), err
		}

		// set header
		if err := pagination.SetPaginationToken(res, linkKey); err != nil {
			return helpers.GatewayErrorResponse(500, "Failed to encode pagination token"), err
		}
	}

	return res, err
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
