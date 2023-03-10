package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	"gitlab.com/deltabyte_/littleurl/api/internal/config"
)

//go:embed templates/*.html
var templates embed.FS

type templateAttributes struct {
	Code        string
	Username    string
	Application string
}

type templateConfig struct {
	FileName string
	Title    string
}

var templateConfigs = map[string]templateConfig{
	"ForgotPassword": {
		FileName: "templates/ForgotPassword.html",
		Title:    "Forgot Password",
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

	// load file
	rawTemplate, err := template.ParseFS(templates, templateConfig.FileName)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return event, err
	}

	// load template
	var template bytes.Buffer
	err = rawTemplate.Execute(&template, &templateAttributes{
		Application: cfg.AppName,
		Code:        event.Request.CodeParameter,
		Username:    event.Request.UsernameParameter,
	})
	if err != nil {
		fmt.Printf("%+v\n", err)
		return event, err
	}

	event.Response.EmailMessage = template.String()
	return event, nil
}

func main() {
	wrappedHander := lumigo.WrapHandler(Handler, &lumigo.Config{})
	lambda.Start(wrappedHander)
}
