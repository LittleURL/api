package templates

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed email/*.html
var emailTemplates embed.FS

func RenderEmail(filename string, attributes interface{}) (string, error) {
	// load file
	rawTemplate, err := template.ParseFS(emailTemplates, fmt.Sprintf("email/%s", filename))
	if err != nil {
		return "", err
	}

	// render
	var template bytes.Buffer
	if err := rawTemplate.Execute(&template, attributes); err != nil {
		return "", err
	}

	return template.String(), nil
}

//go:embed email/cognito/*.html
var cognitoTemplates embed.FS

func RenderCognitoEmail(filename string, attributes interface{}) (string, error) {
	// load file
	rawTemplate, err := template.ParseFS(cognitoTemplates, fmt.Sprintf("email/cognito/%s", filename))
	if err != nil {
		return "", err
	}

	// render
	var template bytes.Buffer
	if err := rawTemplate.Execute(&template, attributes); err != nil {
		return "", err
	}

	return template.String(), nil
}
