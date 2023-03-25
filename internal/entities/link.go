package entities

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"

	av "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gitlab.com/deltabyte_/littleurl/api/internal/timestamp"
)

/**
 * Key
 */
func NewLinkKey(domainId string, uri string) *LinkKey {
	return &LinkKey{
		DomainId: domainId,
		URI:      uri,
	}
}

type LinkKey struct {
	DomainId DomainID `dynamodbav:"domain_id"`
	URI      string   `dynamodbav:"uri"`
}

func (linkKey *LinkKey) MarshalString() (string, error) {
	// encode to gob
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	if err := encoder.Encode(linkKey); err != nil {
		return "", nil
	}

	// encode to base64url
	return base64.RawURLEncoding.EncodeToString(encoded.Bytes()), nil
}

func (linkKey *LinkKey) UnmarshalString(enc string) error {
	// decode base64
	encoded, err := base64.RawURLEncoding.DecodeString(enc)
	if err != nil {
		return err
	}

	// decode gob
	decoded := bytes.NewBuffer(encoded)
	decoder := gob.NewDecoder(decoded)
	if err := decoder.Decode(linkKey); err != nil {
		return nil
	}

	return nil
}

func (linkKey *LinkKey) MarshalDynamoAV() (map[string]ddbTypes.AttributeValue, error) {
	return av.MarshalMap(linkKey)
}

func (linkKey *LinkKey) UnmarshalDynamoAV(item map[string]ddbTypes.AttributeValue) error {
	return av.UnmarshalMap(item, linkKey)
}

/**
 * entity
 */
func NewLink() *Link {
	return &Link{
		CreatedAt: timestamp.Now(),
		UpdatedAt: timestamp.Now(),
	}
}

type Link struct {
	DomainId  DomainID             `json:"domain_id"            dynamodbav:"domain_id"`
	URI       string               `json:"uri"                  dynamodbav:"uri"`
	Target    string               `json:"target"               dynamodbav:"target"`
	CreatedAt timestamp.Timestamp  `json:"created_at"           dynamodbav:"created_at"`
	UpdatedAt timestamp.Timestamp  `json:"updated_at"           dynamodbav:"updated_at"`
	ExpiresAt *timestamp.Timestamp `json:"expires_at,omitempty" dynamodbav:"expires_at"`
}

func (link *Link) MarshalDynamoAV() (map[string]ddbTypes.AttributeValue, error) {
	return av.MarshalMap(link)
}

func (link *Link) UnmarshalDynamoAV(item map[string]ddbTypes.AttributeValue) error {
	return av.UnmarshalMap(item, link)
}

func (link *Link) Key() *LinkKey {
	return &LinkKey{
		DomainId: link.DomainId,
		URI:      link.URI,
	}
}

/**
 * slice of entities
 */
type Links []*Link

func (links *Links) UnmarshalDynamoAV(items []map[string]ddbTypes.AttributeValue) error {
	return av.UnmarshalListOfMaps(items, links)
}
