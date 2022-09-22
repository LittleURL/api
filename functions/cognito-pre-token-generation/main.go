package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"gitlab.com/deltabyte_/littleurl/api/internal/config"
	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
)

func Handler(ctx context.Context, event events.CognitoEventUserPoolsPreTokenGenRequest) (*events.CognitoEventUserPoolsPreTokenGenResponse, error) {
	cfg := config.Load()

	// get user's ID
	userId, exists := event.UserAttributes["sub"]
	if !exists {
		return nil, errors.New("UserID not found")
	}

	// init AWS SDK
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	// get domains
	ddb := dynamodb.NewFromConfig(awsCfg)
	res, err := ddb.Query(ctx, &dynamodb.QueryInput{
		TableName:              &cfg.Tables.DomainUsers,
		IndexName:              aws.String("user-domains"),
		KeyConditionExpression: aws.String("userId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userId},
		},
	})
	if err != nil {
		return nil, err
	}

	// unmarshal
	domainUsers := &entities.DomainUsers{}
	if err := domainUsers.UnmarshalDynamoAV(res.Items); err != nil {
		return nil, err
	}

	// generate claims
	claims := map[string]string{}
	for _, domainUser := range *domainUsers {
		claims[fmt.Sprintf("domain_%s", domainUser.DomainID)] = domainUser.Permission
	}

	// convert response to correct entity
	return &events.CognitoEventUserPoolsPreTokenGenResponse{
		ClaimsOverrideDetails: events.ClaimsOverrideDetails{
			ClaimsToAddOrOverride: claims,
		},
	}, nil
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
