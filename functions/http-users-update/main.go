package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"gitlab.com/deltabyte_/littleurl/api/internal/application"
	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
	"gitlab.com/deltabyte_/littleurl/api/internal/helpers"
)

type UpdateUsersRequest struct {
	RoleName *string `json:"role_name" validate:"required,oneof=admin editor viewer nobody"`
}

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

	// parse body
	reqBody := &UpdateUsersRequest{}
	if err := json.Unmarshal([]byte(event.Body), reqBody); err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to unmarshal body"), err
	}

	// validate request body
	if err := helpers.ValidateRequest(reqBody); err != nil {
		return helpers.GatewayValidationResponse(err), nil
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

	// create entity and marshal
	userRole := &entities.UserRole{
		DomainID: domainId,
		UserID:   userId,
		RoleName: *reqBody.RoleName,
	}
	item, err := userRole.MarshalDynamoAV()
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to marshal userRole"), err
	}

	// write item to dynamo
	_, err = app.DDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &app.Cfg.Tables.UserRoles,
		Item:      item,
		// only allow updates for existing users
		ConditionExpression: aws.String("attribute_exists(domain_id) AND attribute_exists(user_id)"),
	})
	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return helpers.GatewayErrorResponse(400, "You can only update existing roles"), nil
		}

		return helpers.GatewayErrorResponse(500, "Failed to write userRole to database"), err
	}

	return helpers.GatewayJsonResponse(200, nil)
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
