package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
	if err != nil || userRole == nil || !userRole.Role().UsersRead() {
		return helpers.GatewayErrorResponse(403, ""), nil
	}

	// get roles
	userRolesRes, err := app.DDBClient.Query(app.Ctx, &dynamodb.QueryInput{
		TableName:              &app.Cfg.Tables.UserRoles,
		KeyConditionExpression: aws.String("domain_id = :id"),
		ExpressionAttributeValues: map[string]dynamodbTypes.AttributeValue{
			":id": &dynamodbTypes.AttributeValueMemberS{
				Value: domainId,
			},
		},
	})
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to load user roles"), err
	}

	// assume that domain doesn't exist if there are no users, shouldn't be possible given the permission check above
	if len(userRolesRes.Items) < 1 {
		return helpers.GatewayJsonResponse(404, "Domain not found")
	}

	// unmarshal userRoles
	userRoles := entities.UserRoles{}
	if err := userRoles.UnmarshalDynamoAV(userRolesRes.Items); err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to unmarshal user roles"), err
	}

	// load users
	userKeys := []map[string]dynamodbTypes.AttributeValue{}
	for _, userRole := range userRoles {
		userKeys = append(userKeys, map[string]dynamodbTypes.AttributeValue{
			"id": &dynamodbTypes.AttributeValueMemberS{Value: userRole.UserID},
		})
	}
	usersRes, err := app.DDBClient.BatchGetItem(app.Ctx, &dynamodb.BatchGetItemInput{
		RequestItems: map[string]dynamodbTypes.KeysAndAttributes{
			app.Cfg.Tables.Users: {Keys: userKeys},
		},
	})
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to load user details"), err
	}

	// unmarshal users
	users := &entities.Users{}
	if err := users.UnmarshalDynamoAV(usersRes.Responses[app.Cfg.Tables.Users]); err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to unmarshal user details"), err
	}
	for _, user := range *users {
		user.Role = (*userRoles.FindByUserID(user.Id)).RoleName
	}

	return helpers.GatewayJsonResponse(200, users)
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
