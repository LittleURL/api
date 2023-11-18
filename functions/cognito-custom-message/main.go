package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"github.com/littleurl/api/internal/config"
	"github.com/littleurl/api/internal/templates"
)

type templateAttributes struct {
	Subject       string
	Code        string
	Username    string
	Application string
}

type templateConfig struct {
	FileName string
	Subject    string
}

var templateConfigs = map[string]templateConfig{
	"ForgotPassword": {
		FileName: "ForgotPassword.html",
		Subject:    "Forgot Password",
	},
}

func Handler(ctx context.Context, event *events.CognitoEventUserPoolsCustomMessage) (*events.CognitoEventUserPoolsCustomMessage, error) {
	cfg := config.Load()

	// select template config
	name := strings.ReplaceAll(event.TriggerSource, "CustomMessage_", "")
	templateConfig, exists := templateConfigs[name]
	if !exists {
		fmt.Printf("Unknown template: %s", name)
		return event, nil
	}

	// load template
	template, err := templates.RenderCognitoEmail(templateConfig.FileName, &templateAttributes{
		Application: cfg.AppName,
		Code:        event.Request.CodeParameter,
		Subject:       templateConfig.Subject,
		Username:    event.Request.UsernameParameter,
	})
	if err != nil {
		fmt.Printf("%+v\n", err)
		return event, err
	}

	event.Response.EmailMessage = template
	event.Response.EmailSubject = templateConfig.Subject
	return event, nil
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
