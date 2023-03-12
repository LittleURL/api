package entities

import (
	"time"

	av "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"gitlab.com/deltabyte_/littleurl/api/internal/timestamp"
)

type UserInviteID = string

/**
 * entity
 */
type UserInvite struct {
	Id        UserInviteID        `json:"id"         dynamodbav:"id"`
	Email     string              `json:"email"      dynamodbav:"email"`
	Role      string              `json:"role"       dynamodbav:"role"`
	DomainID  DomainID            `json:"domain_id"  dynamodbav:"domain_id"`
	CreatedAt timestamp.Timestamp `json:"created_at" dynamodbav:"created_at"`
	ExpiresAt timestamp.Timestamp `json:"expires_at" dynamodbav:"expires_at"`
}

func (userInvite *UserInvite) MarshalDynamoAV() (map[string]ddbTypes.AttributeValue, error) {
	return av.MarshalMap(userInvite)
}

func (userInvite *UserInvite) UnmarshalDynamoAV(item map[string]ddbTypes.AttributeValue) error {
	return av.UnmarshalMap(item, userInvite)
}

func NewUserInvite() *UserInvite {
	return &UserInvite{
		Id:        uuid.NewString(),
		CreatedAt: timestamp.Now(),
		ExpiresAt: timestamp.Timestamp(time.Now().Add(time.Hour * 24 * 7)),
	}
}

/**
 * slice of entities
 */
type UserInvites []*UserInvite

func (userInvites *UserInvites) UnmarshalDynamoAV(items []map[string]ddbTypes.AttributeValue) error {
	return av.UnmarshalListOfMaps(items, userInvites)
}
