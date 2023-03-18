package entities

import (
	av "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gitlab.com/deltabyte_/littleurl/api/internal/permissions"
	"gitlab.com/deltabyte_/littleurl/api/internal/timestamp"
)

func UserRoleKey(domainId DomainID, userId UserID) map[string]ddbTypes.AttributeValue {
	return map[string]ddbTypes.AttributeValue{
		"domain_id": &ddbTypes.AttributeValueMemberS{Value: domainId},
		"user_id":   &ddbTypes.AttributeValueMemberS{Value: userId},
	}
}

func NewUserRole() *UserRole {
	return &UserRole{
		CreatedAt: timestamp.Now(),
		UpdatedAt: timestamp.Now(),
	}
}

/**
 * entity
 */
type UserRole struct {
	DomainID  DomainID            `json:"domain_id"  dynamodbav:"domain_id"`
	UserID    string              `json:"user_id"    dynamodbav:"user_id"`
	RoleName  string              `json:"role_name"  dynamodbav:"role_name"`
	CreatedAt timestamp.Timestamp `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt timestamp.Timestamp `json:"updated_at" dynamodbav:"updated_at"`
}

func (userRole *UserRole) MarshalDynamoAV() (map[string]ddbTypes.AttributeValue, error) {
	return av.MarshalMap(userRole)
}

func (userRole *UserRole) UnmarshalDynamoAV(item map[string]ddbTypes.AttributeValue) error {
	return av.UnmarshalMap(item, userRole)
}

func (userRole *UserRole) Role() permissions.Role {
	switch userRole.RoleName {
	case permissions.Admin:
		return &permissions.AdminRole{}

	case permissions.Editor:
		return &permissions.EditorRole{}

	case permissions.Viewer:
		return &permissions.ViewerRole{}

	default:
		return &permissions.NobodyRole{}
	}
}

/**
 * slice of entities
 */
type UserRoles []*UserRole

func (userRoles *UserRoles) UnmarshalDynamoAV(items []map[string]ddbTypes.AttributeValue) error {
	return av.UnmarshalListOfMaps(items, userRoles)
}

func (userRoles *UserRoles) FindByUserID(id UserID) *UserRole {
	for _, userRole := range *userRoles {
		if userRole.UserID == id {
			return userRole
		}
	}

	return nil
}

func (userRoles *UserRoles) FindByDomainID(id DomainID) *UserRole {
	for _, userRole := range *userRoles {
		if userRole.DomainID == id {
			return userRole
		}
	}

	return nil
}
