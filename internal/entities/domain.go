package entities

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DomainID = string
type Domain struct {
	Id          DomainID          `json:"user_id"     dynamodbav:"id"`
	Domain      string            `json:"domain"      dynamodbav:"domain"`
	Description string            `json:"description" dynamodbav:"description"`
	Users       map[UserID]string `json:"users" dynamodbav:"users"`
}

func (domain *Domain) MarshalDynamoAV() (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(domain)
}

func (domain *Domain) UnmarshalDynamoAV(av map[string]types.AttributeValue) error {
	return attributevalue.UnmarshalMap(av, domain)
}
