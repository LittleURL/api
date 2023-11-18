package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"github.com/littleurl/api/internal/application"
	"github.com/littleurl/api/internal/entities"
	"github.com/littleurl/api/internal/helpers"
)

func Handler(ctx context.Context, event *events.CognitoEventUserPoolsPreTokenGen) (*events.CognitoEventUserPoolsPreTokenGen, error) {
	// init app
	app, err := application.New(ctx)
	if err != nil {
		fmt.Print("Failed to initialize app")
		panic(err)
	}

	// get user's ID
	userId, exists := event.Request.UserAttributes["sub"]
	if !exists {
		return nil, errors.New("UserID not found")
	}

	// get user details from cognito
	cognitoClient := cognito.NewFromConfig(app.AwsCfg)
	res, err := cognitoClient.ListUsers(ctx, &cognito.ListUsersInput{
		UserPoolId: &app.Cfg.CognitoPoolId,
		Filter:     aws.String(fmt.Sprintf("sub = \"%s\"", userId)),
	})
	if err != nil {
		return nil, err
	}
	if len(res.Users) < 1 {
		fmt.Print("Cognito ListUsers response is empty")
		return nil, fmt.Errorf("failed to fetch user attributes")
	}
	userAttributes := helpers.FlattenCognitoUserAttributes(res.Users[0].Attributes)

	// map to user entity
	fmt.Print(userAttributes)
	user := &entities.User{
		Id:       userId,
		Name:     userAttributes["name"],
		Username: *res.Users[0].Username,
		Nickname: userAttributes["nickname"],
		Picture: userAttributes["picture"],
	}
	ddbUser, err := user.MarshalDynamoAV()
	if err != nil {
		return nil, err
	}

	// get domains
	_, err = app.DDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &app.Cfg.Tables.Users,
		Item:      ddbUser,
	})
	if err != nil {
		return nil, err
	}

	return event, nil
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
