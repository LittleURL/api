package entities

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserID = string

type Users []User

func (users *Users) UnmarshalDynamoAV(items []map[string]types.AttributeValue) error {
	return attributevalue.UnmarshalListOfMaps(items, users)
}

type User struct {
	Id       UserID `json:"id"       dynamodbav:"id"`
	Name     string `json:"name"     dynamodbav:"name"`
	Username string `json:"username" dynamodbav:"username"`
}

func (user *User) MarshalDynamoAV() (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(user)
}

func (user *User) UnmarshalDynamoAV(item map[string]types.AttributeValue) error {
	return attributevalue.UnmarshalMap(item, user)
}
