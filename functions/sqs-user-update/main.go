package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/auth0/go-auth0/management"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/deltabyte/littleurl-api/internal/config"
	"github.com/deltabyte/littleurl-api/internal/entities"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
)

func Handler(ctx context.Context, event events.SQSEvent) error {
	// load config
	cfg := config.Load()

	// init auth0 sdk
	auth0Client, err := management.New(
		cfg.Auth0.Domain,
		management.WithClientCredentials(cfg.Auth0.ClientID, cfg.Auth0.ClientSecret),
	)
	if err != nil {
		return err
	}

	for _, message := range event.Records {
		userId, ok := message.MessageAttributes["user_id"]
		if !ok {
			return fmt.Errorf("failed to read UserID")
		}

		var user entities.User
		if err := json.Unmarshal([]byte(message.Body), &user); err != nil {
			return err
		}

		auth0User, err := user.ToAuth0User()
		if err != nil {
			return err
		}

		fmt.Printf("Updating user: %s", *userId.StringValue)
		if err := auth0Client.User.Update(*userId.StringValue, auth0User); err != nil {
			return err
		}

		// fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
	}

	// success
	return nil
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
