package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"github.com/littleurl/api/internal/application"
	"github.com/littleurl/api/internal/entities"
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
	invitesRepo := repositories.NewUserInvitesRepository(app)
	rolesRepo := repositories.NewUserRolesRepository(app)

	// extract the user's ID
	currentUserId, exists := event.RequestContext.Authorizer.JWT.Claims["sub"]
	if !exists {
		return helpers.GatewayErrorResponse(401, "UserID not found in auth token"), nil
	}

	// validate path params
	inviteId, exists := event.PathParameters["inviteId"]
	if !exists {
		return helpers.GatewayErrorResponse(400, "Missing `inviteId` path parameter"), nil
	}

	// get user details from cognito
	cognitoClient := cognito.NewFromConfig(app.AwsCfg)
	listUsersRes, err := cognitoClient.ListUsers(ctx, &cognito.ListUsersInput{
		UserPoolId: &app.Cfg.CognitoPoolId,
		Filter:     aws.String(fmt.Sprintf("sub = \"%s\"", currentUserId)),
	})
	if err != nil {
		return nil, err
	}
	if len(listUsersRes.Users) < 1 {
		return helpers.GatewayErrorResponse(500, ""), fmt.Errorf("Cognito ListUsers response is empty")
	}
	userAttributes := helpers.FlattenCognitoUserAttributes(listUsersRes.Users[0].Attributes)

	// check user can confirmed their email
	if listUsersRes.Users[0].UserStatus != "CONFIRMED" {
		return helpers.GatewayErrorResponse(401, "User is not confirmed"), nil
	}

	// get invite
	invite, reqErr := invitesRepo.Find(ctx, inviteId)
	if reqErr != nil {
		return reqErr.GatewayResponse(), reqErr.Error()
	}

	// validate invite
	if invite.Expired() || invite.Email != userAttributes["email"] {
		return helpers.GatewayErrorResponse(403, "Invalid invite"), nil
	}

	// create role
	userRole := entities.NewUserRole()
	userRole.DomainID = invite.DomainID
	userRole.UserID = currentUserId
	userRole.RoleName = invite.Role
	if reqErr := rolesRepo.Create(ctx, userRole); reqErr != nil {
		return reqErr.GatewayResponse(), reqErr.Error()
	}

	// delete invites
	reqErr = invitesRepo.Cleanup(ctx, invite)
	if reqErr != nil {
		return reqErr.GatewayResponse(), reqErr.Error()
	}

	return helpers.GatewayJsonResponse(200, nil)
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
