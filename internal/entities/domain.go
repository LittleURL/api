package entities

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Domain struct {
	Id          string `json:"user_id"     dynamodbav:"id"`
	Domain      string `json:"domain"      dynamodbav:"domain"`
	Description string `json:"description" dynamodbav:"description"`
	Favicon     []byte `json:"favicon"     dynamodbav:"favicon"`
}

func (domain *Domain) MarshalDynamoAV() (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(domain)
}

func (domain *Domain) UnmarshalDynamoAV(av map[string]types.AttributeValue) error {
	return attributevalue.UnmarshalMap(av, domain)
}
