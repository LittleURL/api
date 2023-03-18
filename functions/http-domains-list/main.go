package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

	// init AWS SDK
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	ddb := dynamodb.NewFromConfig(awsCfg)

	// get domain IDs that the current user has access to
	userRolesRes, err := ddb.Query(ctx, &dynamodb.QueryInput{
		TableName:              &app.Cfg.Tables.UserRoles,
		IndexName:              aws.String("user-domains"),
		KeyConditionExpression: aws.String("user_id = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{
				Value: userId,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	userRoles := entities.UserRoles{}
	if err := userRoles.UnmarshalDynamoAV(userRolesRes.Items); err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to parse domain permissions"), err
	}

	// filter roles to only those with real permissions
	allowedRoles := entities.UserRoles{}
	for _, userRole := range userRoles {
		if userRole.Role().DomainRead() {
			allowedRoles = append(allowedRoles, userRole)
		}
	}

	// shortcut rest of function if the user has zero roles
	if len(allowedRoles) < 1 {
		return helpers.GatewayJsonResponse(200, make([]string, 0))
	}

	// extract domain IDs
	domainKeys := []map[string]types.AttributeValue{}
	for _, userRole := range allowedRoles {
		domainKeys = append(domainKeys, map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: userRole.DomainID,
			},
		})
	}

	// get domains
	res, err := ddb.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			app.Cfg.Tables.Domains: {Keys: domainKeys},
		},
	})
	if err != nil {
		panic(err)
	}

	// convert response to correct entity
	domains := []*entities.Domain{}
	for _, domain := range res.Responses[app.Cfg.Tables.Domains] {
		newDomain := &entities.Domain{}
		if err := newDomain.UnmarshalDynamoAV(domain); err != nil {
			panic(err)
		}

		// add the current user's role
		newDomain.UserRole = (*userRoles.FindByDomainID(newDomain.Id)).RoleName

		domains = append(domains, newDomain)
	}

	// send data to user
	return helpers.GatewayJsonResponse(200, domains)
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
