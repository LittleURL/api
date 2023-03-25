package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jinzhu/copier"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"gitlab.com/deltabyte_/littleurl/api/internal/application"
	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
	"gitlab.com/deltabyte_/littleurl/api/internal/helpers"
)

type CreateLinkRequest struct {
	URI       string  `json:"uri"                  validate:"required,max=1024,uri,startsnotwith=/_littleurl"`
	Target    string  `json:"target"               validate:"required,max=2048,url"`
	ExpiresAt *string `json:"expires_at,omitempty" validate:omitempty,datetime"`
}

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

	// parse body
	reqBody := &CreateLinkRequest{}
	if err := json.Unmarshal([]byte(event.Body), reqBody); err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to unmarshal body"), err
	}

	// validate request body
	if err := helpers.ValidateRequest(reqBody); err != nil {
		return helpers.GatewayValidationResponse(err), nil
	}

	// check permissions
	userRole, err := helpers.FindUserRole(app, domainId, userId)
	if err != nil || userRole == nil || !userRole.Role().LinksWrite() {
		return helpers.GatewayErrorResponse(403, ""), err
	}

	// create entity
	link := entities.NewLink()
	link.DomainId = domainId
	if err := copier.Copy(link, reqBody); err != nil {
		return helpers.GatewayErrorResponse(500, ""), err
	}

	// marshal item
	linkItem, err := link.MarshalDynamoAV()
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to marshal link"), err
	}

	// get roles
	_, err = app.DDBClient.PutItem(app.Ctx, &dynamodb.PutItemInput{
		TableName: &app.Cfg.Tables.Links,
		Item:      linkItem,
		// condition expression not required due to ddb pk/sk uniqueness
	})
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to save link"), err
	}

	return helpers.GatewayJsonResponse(201, link)
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
