package entities

import (
	av "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gitlab.com/deltabyte_/littleurl/api/internal/timestamp"
)

func LinkKey(domainId DomainID, uri string) map[string]ddbTypes.AttributeValue {
	return map[string]ddbTypes.AttributeValue{
		"domain_id": &ddbTypes.AttributeValueMemberS{Value: domainId},
		"uri":   &ddbTypes.AttributeValueMemberS{Value: uri},
	}
}

/**
 * entity
 */
type Link struct {
	DomainId  DomainID            `json:"domain_id"  dynamodbav:"domain_id"`
	URI       string              `json:"uri"        dynamodbav:"uri"`
	Target    string              `json:"target"     dynamodbav:"target"`
	CreatedAt timestamp.Timestamp `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt timestamp.Timestamp `json:"updated_at" dynamodbav:"updated_at"`
}

func NewLink() *Link {
	return &Link{
		CreatedAt: timestamp.Now(),
		UpdatedAt: timestamp.Now(),
	}
}

func (link *Link) MarshalDynamoAV() (map[string]ddbTypes.AttributeValue, error) {
	return av.MarshalMap(link)
}

func (link *Link) UnmarshalDynamoAV(item map[string]ddbTypes.AttributeValue) error {
	return av.UnmarshalMap(item, link)
}

/**
 * slice of entities
 */
type Links []*Link

func (links *Links) UnmarshalDynamoAV(items []map[string]ddbTypes.AttributeValue) error {
	return av.UnmarshalListOfMaps(items, links)
}
