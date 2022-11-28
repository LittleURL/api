package entities

import (
	av "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserID = string

/**
 * entity
 */
type User struct {
	Id       UserID `json:"id"             dynamodbav:"id"`
	Name     string `json:"name"           dynamodbav:"name"`
	Username string `json:"username"       dynamodbav:"username"`
	Picture  string `json:"picture"        dynamodbav:"picture"`
	Role     string `json:"role,omitempty" dynamodbav:"-"` // not stored in same table as users
}

func (user *User) MarshalDynamoAV() (map[string]ddbTypes.AttributeValue, error) {
	return av.MarshalMap(user)
}

func (user *User) UnmarshalDynamoAV(item map[string]ddbTypes.AttributeValue) error {
	return av.UnmarshalMap(item, user)
}

/**
 * slice of entities
 */
type Users []*User

func (users *Users) UnmarshalDynamoAV(items []map[string]ddbTypes.AttributeValue) error {
	return av.UnmarshalListOfMaps(items, users)
}
