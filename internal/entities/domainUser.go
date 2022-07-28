package entities

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DomainUser struct {
	DomainID   DomainID `json:"domain_id"  dynamodbav:"domain_id"`
	UserID     string   `json:"user_id"    dynamodbav:"user_id"`
	Permission string   `json:"permission" dynamodbav:"permission"`
}

func (domainUser *DomainUser) MarshalDynamoAV() (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(domainUser)
}

func (domainUser *DomainUser) UnmarshalDynamoAV(item map[string]types.AttributeValue) error {
	return attributevalue.UnmarshalMap(item, domainUser)
}

type DomainUsers []DomainUser
func (domainUsers *DomainUsers) UnmarshalDynamoAV(items []map[string]types.AttributeValue) error {
	return attributevalue.UnmarshalListOfMaps(items, domainUsers)
}
