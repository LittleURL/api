package entities

import (
	av "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/littleurl/api/internal/timestamp"
)

type DomainID = string

/**
 * entity
 */
type Domain struct {
	// core data
	Id            DomainID `json:"id"                       dynamodbav:"id"`
	Domain        string   `json:"domain"                   dynamodbav:"domain"`
	Description   *string  `json:"description,omitempty"    dynamodbav:"description,omitempty"`
	DefaultTarget *string  `json:"default_target,omitempty" dynamodbav:"default_target,omitempty"`

	// engine config
	Engine EngineConfig `json:"engine" dynamodbav:"engine"`

	// meta
	CreatedAt timestamp.Timestamp `json:"created_at"               dynamodbav:"created_at"`
	UpdatedAt timestamp.Timestamp `json:"updated_at"               dynamodbav:"updated_at"`
	UserRole  string              `json:"user_role,omitempty"      dynamodbav:"-"` // not stored in dynamo
}

func NewDomain() *Domain {
	return &Domain{
		Id:        uuid.NewString(),
		CreatedAt: timestamp.Now(),
		UpdatedAt: timestamp.Now(),
	}
}

func (domain *Domain) MarshalDynamoAV() (map[string]ddbTypes.AttributeValue, error) {
	return av.MarshalMap(domain)
}

func (domain *Domain) UnmarshalDynamoAV(item map[string]ddbTypes.AttributeValue) error {
	return av.UnmarshalMap(item, domain)
}
