package entities

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gitlab.com/deltabyte_/littleurl/api/internal/permissions"
)

/**
 * entity
 */
type UserRole struct {
	DomainID DomainID `json:"domain_id" dynamodbav:"domain_id"`
	UserID   string   `json:"user_id"   dynamodbav:"user_id"`
	RoleName string   `json:"role_name" dynamodbav:"role_name"`
}

func (userRole *UserRole) MarshalDynamoAV() (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(userRole)
}

func (userRole *UserRole) UnmarshalDynamoAV(item map[string]types.AttributeValue) error {
	return attributevalue.UnmarshalMap(item, userRole)
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

func (userRoles *UserRoles) UnmarshalDynamoAV(items []map[string]types.AttributeValue) error {
	return attributevalue.UnmarshalListOfMaps(items, userRoles)
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