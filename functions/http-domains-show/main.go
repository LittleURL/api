package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"github.com/littleurl/api/internal/application"
	"github.com/littleurl/api/internal/helpers"
	"github.com/littleurl/api/internal/repositories"
)

func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	// init app
	app, err := application.New(ctx)
	if err != nil {
		fmt.Print("Failed to initialize app")
		panic(err)
	}
	domainsRepo := repositories.NewDomainsRepository(app)

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
		return helpers.GatewayErrorResponse(403, ""), err
	}

	// get domain
	domain, reqErr := domainsRepo.Find(ctx, domainId)
	if reqErr != nil {
		return reqErr.GatewayResponse(), reqErr.Err
	}

	// set current user's role
	domain.UserRole = userRole.RoleName

	return helpers.GatewayJsonResponse(200, domain)
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
