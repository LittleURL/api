package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/ses/types"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"gitlab.com/deltabyte_/littleurl/api/internal/application"
	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
	"gitlab.com/deltabyte_/littleurl/api/internal/helpers"
	"gitlab.com/deltabyte_/littleurl/api/internal/repositories"
	"gitlab.com/deltabyte_/littleurl/api/internal/templates"
)

type InviteUserRequest struct {
	Email *string `json:"email" validate:"required,email"`
	Role  *string `json:"role" validate:"required,oneof=admin editor viewer nobody"`
}

func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	// init app
	app, err := application.New(ctx)
	if err != nil {
		fmt.Print("Failed to initialize app")
		panic(err)
	}
	domainsRepo := repositories.NewDomainsRepository(app)

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
	reqBody := &InviteUserRequest{}
	if err := json.Unmarshal([]byte(event.Body), reqBody); err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to unmarshal body"), err
	}

	// validate request body
	if err := helpers.ValidateRequest(reqBody); err != nil {
		return helpers.GatewayValidationResponse(err), nil
	}

	// check permissions
	currentUserRole, err := helpers.FindUserRole(app, domainId, currentUserId)
	if err != nil || currentUserRole == nil || !currentUserRole.Role().UsersWrite() {
		return helpers.GatewayErrorResponse(403, ""), err
	}

	// get domain
	domain, reqErr := domainsRepo.Find(ctx, domainId)
	if reqErr != nil {
		return reqErr.GatewayResponse(), reqErr.Err
	}

	// create entity
	userInvite := entities.NewUserInvite()
	userInvite.Email = *reqBody.Email
	userInvite.Role = *reqBody.Role
	userInvite.DomainID = domainId

	// marshal entity
	item, err := userInvite.MarshalDynamoAV()
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to marshal userInvite"), err
	}

	// write item to dynamo
	_, err = app.DDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           &app.Cfg.Tables.UserInvites,
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	})
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to write userInvite to database"), err
	}

	// render email template
	emailHtml, err := templates.RenderEmail("InviteUser.html", map[string]string{
		"Application": app.Cfg.AppName,
		"Domain":      domain.Domain,
		"Link":        fmt.Sprintf("https://%s/invite/%s", app.Cfg.DashboardDomain, userInvite.ID),
		"Expiry":      fmt.Sprintf("%.0f days", userInvite.ExpiresAt.Until().Hours() / 24),
	})
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to render email"), err
	}

	// send email
	sesClient := ses.NewFromConfig(app.AwsCfg)
	_, err = sesClient.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &sesTypes.Destination{
			ToAddresses: []string{userInvite.Email},
		},
		Message: &sesTypes.Message{
			Body: &sesTypes.Body{
				Html: &sesTypes.Content{Data: &emailHtml},
				// TODO: text based body
			},
			Subject: &sesTypes.Content{
				Data: aws.String(fmt.Sprintf("%s Invitation", app.Cfg.AppName)),
			},
		},
		Source: &app.Cfg.EmailFrom,
	})
	if err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to send email"), err
	}

	return helpers.GatewayJsonResponse(200, nil)
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
