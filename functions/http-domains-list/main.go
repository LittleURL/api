package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"gitlab.com/deltabyte_/littleurl/api/internal/config"
	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
)

func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	cfg := config.Load()

	// check allowed domains
	claims := event.RequestContext.Authorizer.JWT.Claims
	domainKeys := []map[string]types.AttributeValue{}
	for k := range claims {
		if strings.HasPrefix(k, "domain_") {
			id := strings.Replace(k, "domain_", "", 1)
			domainKeys = append(domainKeys, map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: id},
			})
		}
	}

	// init AWS SDK
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	// get domains
	ddb := dynamodb.NewFromConfig(awsCfg)
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
	body, err := json.Marshal(domains)
	if err != nil {
		panic(err)
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
