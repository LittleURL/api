package entities

import (
	"time"

	av "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"gitlab.com/deltabyte_/littleurl/api/internal/timestamp"
)

func NewUserInvite() *UserInvite {
	return &UserInvite{
		ID:        uuid.NewString(),
		CreatedAt: timestamp.Now(),
		ExpiresAt: timestamp.Timestamp(time.Now().Add(time.Hour * 24 * 7)),
	}
}

/**
 * entity
 */
type UserInviteID = string
type UserInvite struct {
	ID        UserInviteID        `json:"id"         dynamodbav:"id"`
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

func (userInvite *UserInvite) Expired() bool {
	return userInvite.ExpiresAt.Until().Nanoseconds() == 0
}


/**
 * slice of entities
 */
type UserInvites []*UserInvite

func (userInvites *UserInvites) UnmarshalDynamoAV(items []map[string]ddbTypes.AttributeValue) error {
	return av.UnmarshalListOfMaps(items, userInvites)
}
