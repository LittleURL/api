package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitoTypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"gitlab.com/deltabyte_/littleurl/api/internal/application"
	"gitlab.com/deltabyte_/littleurl/api/internal/helpers"
)

func gravatarUrl(email string) string {
	clean := strings.ToLower(strings.TrimSpace(email))
	hash := md5.Sum([]byte(clean))
	emailHash := hex.EncodeToString(hash[:])
	return "https://www.gravatar.com/avatar/" + emailHash
}


func Handler(ctx context.Context, event *events.CognitoEventUserPoolsPostAuthentication) (*events.CognitoEventUserPoolsPostAuthentication, error) {
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
	listUsersRes, err := cognitoClient.ListUsers(ctx, &cognito.ListUsersInput{
		UserPoolId: &app.Cfg.CognitoPoolId,
		Filter:     aws.String(fmt.Sprintf("sub = \"%s\"", userId)),
	})
	if err != nil {
		return nil, err
	}
	if len(listUsersRes.Users) < 1 {
		fmt.Print("Cognito ListUsers response is empty")
		return nil, fmt.Errorf("failed to fetch user attributes")
	}
	userAttributes := helpers.FlattenCognitoUserAttributes(listUsersRes.Users[0].Attributes)

	// gravatar
	updatedAttributes := []cognitoTypes.AttributeType{}
	userEmail, emailExists := userAttributes["email"]
	if emailExists {
		updatedAttributes = append(updatedAttributes, cognitoTypes.AttributeType{
			Name: aws.String("picture"),
			Value: aws.String(gravatarUrl(userEmail)),
		})
	}
	
	// update the user's attributes
	_, err = cognitoClient.AdminUpdateUserAttributes(ctx, &cognito.AdminUpdateUserAttributesInput{
		UserPoolId: &app.Cfg.CognitoPoolId,
		Username: listUsersRes.Users[0].Username,
		UserAttributes: updatedAttributes,
	})
	if err != nil {
		return event, err
	}

	return event, nil
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
