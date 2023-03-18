package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"gitlab.com/deltabyte_/littleurl/api/internal/application"
	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
	"gitlab.com/deltabyte_/littleurl/api/internal/helpers"
	"gitlab.com/deltabyte_/littleurl/api/internal/repositories"
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
	rolesRepo := repositories.NewUserRolesRepository(app)

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
	if reqErr := rolesRepo.Update(ctx, userRole); reqErr != nil {
		return reqErr.GatewayResponse(), reqErr.Err
	}

	return helpers.GatewayJsonResponse(200, nil)
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
