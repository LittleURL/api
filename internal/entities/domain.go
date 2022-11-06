package entities

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gitlab.com/deltabyte_/littleurl/api/internal/application"
)

type DomainID = string

/**
 * slice of entities
 */
type Domain struct {
	Id          DomainID `json:"id"              dynamodbav:"id"`
	Domain      string   `json:"domain"          dynamodbav:"domain"`
	Description string   `json:"description"     dynamodbav:"description"`
	Users       *Users   `json:"users,omitempty" dynamodbav:"-"` // not stored in this dynamodb table
}

func (domain *Domain) MarshalDynamoAV() (map[string]dynamodbTypes.AttributeValue, error) {
	return attributevalue.MarshalMap(domain)
}

func (domain *Domain) UnmarshalDynamoAV(av map[string]dynamodbTypes.AttributeValue) error {
	return attributevalue.UnmarshalMap(av, domain)
}

func (domain *Domain) LoadUsers(app *application.Application) error {
	// get permissions
	userRolesRes, err := app.DDBClient.Query(app.Ctx, &dynamodb.QueryInput{
		TableName:              &app.Cfg.Tables.UserRoles,
		KeyConditionExpression: aws.String("domain_id = :id"),
		ExpressionAttributeValues: map[string]dynamodbTypes.AttributeValue{
			":id": &dynamodbTypes.AttributeValueMemberS{
				Value: domain.Id,
			},
		},
	})
	if err != nil {
		return err
	}

	// unmarshal userRoles
	userRoles := UserRoles{}
	if err := userRoles.UnmarshalDynamoAV(userRolesRes.Items); err != nil {
		return err
	}

	// load users
	userKeys := []map[string]dynamodbTypes.AttributeValue{}
	for _, userRole := range userRoles {
		userKeys = append(userKeys, map[string]dynamodbTypes.AttributeValue{
			"id": &dynamodbTypes.AttributeValueMemberS{Value: userRole.UserID},
		})
	}
	usersRes, err := app.DDBClient.BatchGetItem(app.Ctx, &dynamodb.BatchGetItemInput{
		RequestItems: map[string]dynamodbTypes.KeysAndAttributes{
			app.Cfg.Tables.Users: {Keys: userKeys},
		},
	})
	if err != nil {
		return err
	}

	// unmarshal users
	users := &Users{}
	if err := users.UnmarshalDynamoAV(usersRes.Responses[app.Cfg.Tables.Users]); err != nil {
		return err
	}
	for _, user := range *users {
		user.Role = (*userRoles.FindByUserID(user.Id)).RoleName
	}

	// append to domain
	domain.Users = users
	return nil
}
