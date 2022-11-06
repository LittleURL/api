package entities

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DomainID = string

/**
 * entity
 */
type Domain struct {
	Id          DomainID `json:"id"                  dynamodbav:"id"`
	Domain      string   `json:"domain"              dynamodbav:"domain"`
	Description string   `json:"description"         dynamodbav:"description"`
	UserRole    string   `json:"user_role,omitempty" dynamodbav:"-"` // not stored in dynamo
}

func (domain *Domain) MarshalDynamoAV() (map[string]dynamodbTypes.AttributeValue, error) {
	return attributevalue.MarshalMap(domain)
}

func (domain *Domain) UnmarshalDynamoAV(av map[string]dynamodbTypes.AttributeValue) error {
	return attributevalue.UnmarshalMap(av, domain)
}
