package main

import (
	"context"
	"encoding/base64"
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
	userId, exists := event.RequestContext.Authorizer.JWT.Claims["sub"]
	if !exists {
		return helpers.GatewayErrorResponse(401, "UserID not found in auth token"), nil
	}

	// validate path params
	domainId, exists := event.PathParameters["domainId"]
	if !exists {
		return helpers.GatewayErrorResponse(400, "Missing `domainId` path parameter"), nil
	}
	uriParam, exists := event.PathParameters["uri"]
	if !exists {
		return helpers.GatewayErrorResponse(400, "Missing `uri` path parameter"), nil
	}

	// decode URI
	fmt.Print(uriParam)
	decodedUri, err := base64.RawURLEncoding.DecodeString(uriParam)
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to decode URI"), err
	}
	uri := string(decodedUri)

	// check permissions
	userRole, err := helpers.FindUserRole(app, domainId, userId)
	if err != nil || userRole == nil || !userRole.Role().LinksWrite() {
		return helpers.GatewayErrorResponse(403, ""), err
	}

	// get roles
	_, err = app.DDBClient.DeleteItem(app.Ctx, &dynamodb.DeleteItemInput{
		TableName: &app.Cfg.Tables.Links,
		Key:       entities.LinkKey(domainId, uri),
	})
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to delete link"), err
	}

	return helpers.GatewayJsonResponse(200, nil)
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
