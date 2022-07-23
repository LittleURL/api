package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/deltabyte/littleurl-api/internal/config"
	"github.com/deltabyte/littleurl-api/internal/helpers"
	"github.com/deltabyte/littleurl-api/internal/permissions"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
)

type UpdateDomainRequest struct {
	Description *string `json:"description" validate:"max=255"`
}

func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
	cfg := config.Load()

	// validate path params
	domainId, exists := event.PathParameters["domainId"]
	if !exists {
		return helpers.GatewayErrorResponse(400, "Missing `domainId` path parameter"), nil
	}

	// check permissions
	claims := event.RequestContext.Authorizer.JWT.Claims
	userId, userIdExists := claims["sub"]
	roles := permissions.ParseClaims(claims)
	domainRole, roleExists := roles[domainId]
	if !userIdExists || !roleExists || !domainRole.DomainWrite(domainId) {
		return helpers.GatewayErrorResponse(403, ""), nil
	}

	// parse body
	reqBody := &UpdateDomainRequest{}
	if err := json.Unmarshal([]byte(event.Body), reqBody); err != nil {
		return helpers.GatewayErrorResponse(500, "Failed to unmarshal body"), err
	}

	// validate request body
	if err := helpers.ValidateRequest(reqBody); err != nil {
		return helpers.GatewayValidationResponse(err), nil
	}

	// init AWS SDK
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	// build expression
	builder := expression.NewBuilder()
	builder = helpers.ChangeTrackingExpression(builder, userId)

	// update description
	if reqBody.Description != nil {
		builder = builder.WithUpdate(expression.Set(
			expression.Name("description"),
			expression.Value(reqBody.Description),
		))
	}

	// compile expression
	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}

	// get domains
	ddb := dynamodb.NewFromConfig(awsCfg)
	_, err = ddb.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &cfg.Tables.Domains,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: domainId},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
	})
	if err != nil {
		panic(err)
	}

	return &events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       "",
	}, nil
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
