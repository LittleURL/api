package entities

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/auth0/go-auth0/management"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

type User struct {
	Id          string      `required:"true" json:"user_id"`
	AppMetadata AppMetadata `json:"app_metadata,omitempty"`
}

func (user *User) ToAuth0User() (*management.User, error) {
	userJson, err := json.Marshal(*user)
	if err != nil {
		return nil, err
	}

	auth0User := &management.User{}
	if err := auth0User.UnmarshalJSON(userJson); err != nil {
		return nil, err
	}

	// auth0 is stupid and I hate their api
	if !strings.HasPrefix(*auth0User.ID, "auth0|") {
		prefixed := fmt.Sprintf("auth0|%s", *auth0User.ID)
		auth0User.ID = &prefixed
	}

	return auth0User, nil
}

func (user User) MarshalSQSBatchEvent() (*sqs.SendMessageBatchRequestEntry, error) {
	userJson, err := json.Marshal(&user)
	if err != nil {
		return nil, err
	}
	userBase64 := base64.StdEncoding.EncodeToString(userJson)

	sqsEvent := &sqs.SendMessageBatchRequestEntry{
		Id:             aws.String(uuid.NewString()),
		MessageGroupId: aws.String(user.Id),
		MessageBody:    aws.String(userBase64),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"user_id": {
				DataType:    aws.String("String"),
				StringValue: aws.String(user.Id),
			},
		},
	}

	return sqsEvent, nil
}

func (user *User) UnmarshalSQSMessageBody(b64 string) error {
	rawJson, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}

	return json.Unmarshal(rawJson, user)
}
