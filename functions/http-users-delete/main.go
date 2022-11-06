package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
	currentUserId, exists := event.RequestContext.Authorizer.JWT.Claims["sub"]
	if !exists {
		return helpers.GatewayErrorResponse(401, "UserID not found in auth token"), nil
	}

	// validate path params
	domainId, exists := event.PathParameters["domainId"]
	if !exists {
		return helpers.GatewayErrorResponse(400, "Missing `domainId` path parameter"), nil
	}
	userId, exists := event.PathParameters["userId"]
	if !exists {
		return helpers.GatewayErrorResponse(400, "Missing `userId` path parameter"), nil
	}

	// prevent changing of own perms
	if userId == currentUserId {
		return helpers.GatewayErrorResponse(400, "Cannot change own role"), nil
	}

	// check permissions
	currentUserRole, err := helpers.FindUserRole(app, domainId, currentUserId)
	if err != nil || currentUserRole == nil || !currentUserRole.Role().UsersWrite() {
		return helpers.GatewayErrorResponse(403, ""), err
	}

	// write item to dynamo
	_, err = app.DDBClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &app.Cfg.Tables.UserRoles,
		Key:       entities.UserRoleKey(domainId, userId),
	})
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to delete userRole"), err
	}

	return helpers.GatewayJsonResponse(200, nil)
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
