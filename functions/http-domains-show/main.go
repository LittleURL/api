package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"gitlab.com/deltabyte_/littleurl/api/internal/application"
	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
	"gitlab.com/deltabyte_/littleurl/api/internal/helpers"
)

func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
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
	if err != nil || userRole == nil || !userRole.Role().DomainRead() {
		return helpers.GatewayErrorResponse(403, ""), nil
	}

	// get domain
	res, err := app.DDBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &app.Cfg.Tables.Domains,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: domainId,
			},
		},
	})
	if err != nil || res == nil || res.Item == nil {
		return helpers.GatewayErrorResponse(404, "Domain not found"), err
	}
	domain := &entities.Domain{}
	if err := domain.UnmarshalDynamoAV(res.Item); err != nil {
		panic(err) // TODO: gracefully handle unmarshalling error
	}

	// set current user's role
	domain.UserRole = userRole.RoleName

	return helpers.GatewayJsonResponse(200, domain)
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
