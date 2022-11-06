package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"gitlab.com/deltabyte_/littleurl/api/internal/application"
	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
	"gitlab.com/deltabyte_/littleurl/api/internal/helpers"
	"gitlab.com/deltabyte_/littleurl/api/internal/permissions"
	"golang.org/x/exp/slices"
)

const DeleteRole string = "DELETE"

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

	// parse body
	reqBody := &map[entities.UserID]string{}
	if err := json.Unmarshal([]byte(event.Body), reqBody); err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to unmarshal body"), err
	}

	// validate request body
	for userId, roleName := range *reqBody {
		// prevent user from modifying themselves
		if userId == currentUserId {
			return helpers.GatewayErrorResponse(400, "You cannot modify your own role"), nil
		}

		// check roleName is
		if roleName != DeleteRole || !slices.Contains(permissions.RoleNames, roleName) {
			return helpers.GatewayErrorResponse(400, fmt.Sprintf("Invalid role, must be one of: [%s]",
				strings.Join(permissions.RoleNames, ", "),
			)), nil
		}
	}

	// check permissions
	userRole, err := helpers.FindUserRole(app, domainId, currentUserId)
	if err != nil || userRole == nil || !userRole.Role().DomainWrite() {
		return helpers.GatewayErrorResponse(403, ""), nil
	}

	// marshal data
	writeRequests := []dynamodbTypes.WriteRequest{}
	for userId, roleName := range *reqBody {
		// delete userRole
		if roleName == DeleteRole {
			writeRequests = append(writeRequests, dynamodbTypes.WriteRequest{
				DeleteRequest: &dynamodbTypes.DeleteRequest{
					Key: entities.UserRoleKey(domainId, userId),
				},
			})

			continue
		}

		// update userRole
		userRole := &entities.UserRole{
			DomainID: domainId,
			UserID:   userId,
			RoleName: roleName,
		}

		// marshal to dynamo
		item, err := userRole.MarshalDynamoAV()
		if err != nil {
			return helpers.GatewayErrorResponse(500, "Failed to marshal userRole"), err
		}

		writeRequests = append(writeRequests, dynamodbTypes.WriteRequest{
			PutRequest: &dynamodbTypes.PutRequest{Item: item},
		})
	}

	// update items
	_, err = app.DDBClient.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]dynamodbTypes.WriteRequest{
			app.Cfg.Tables.UserRoles: writeRequests,
		},
	})
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to update userRoles table"), err
	}

	return helpers.GatewayJsonResponse(200, nil)
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
