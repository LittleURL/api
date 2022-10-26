package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"gitlab.com/deltabyte_/littleurl/api/internal/config"
	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
	"gitlab.com/deltabyte_/littleurl/api/internal/helpers"
)

func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	cfg := config.Load()

	// extract the user's ID
	userId, exists := event.RequestContext.Authorizer.JWT.Claims["kid"]
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
	userDomainsRes, err := ddb.Query(ctx, &dynamodb.QueryInput{
		TableName:        &cfg.Tables.DomainUsers,
		IndexName:        aws.String("user-domains"),
		FilterExpression: aws.String("user_id = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{
				Value: userId,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	domainUsers := entities.DomainUsers{}
	if err := domainUsers.UnmarshalDynamoAV(userDomainsRes.Items); err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to parse domain permissions"), err
	}

	// extract domain IDs
	domainKeys := []map[string]types.AttributeValue{}
	for _, domainUser := range domainUsers {
		if domainUser.Role().DomainRead() {
			domainKeys = append(domainKeys, map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{
					Value: domainUser.DomainID,
				},
			})
		}
	}

	// get domains
	res, err := ddb.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			cfg.Tables.Domains: {Keys: domainKeys},
		},
	})
	if err != nil {
		panic(err)
	}

	// convert response to correct entity
	domains := []*entities.Domain{}
	for _, domain := range res.Responses[cfg.Tables.Domains] {
		newDomain := &entities.Domain{}
		if err := newDomain.UnmarshalDynamoAV(domain); err != nil {
			panic(err)
		}
		domains = append(domains, newDomain)
	}

	// send data to user
	return helpers.GatewayJsonResponse(200, domains)
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
