package repositories

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/littleurl/api/internal/application"
	"github.com/littleurl/api/internal/entities"
)

type UserRolesRepository struct {
	TableName *string
	DDBClient *dynamodb.Client
}

func NewUserRolesRepository(app *application.Application) *UserRolesRepository {
	return &UserRolesRepository{
		TableName: aws.String(app.Cfg.Tables.UserRoles),
		DDBClient: app.DDBClient,
	}
}

func (repo *UserRolesRepository) Create(ctx context.Context, userRole *entities.UserRole) *application.RequestError {
	// marshal entity
	item, err := userRole.MarshalDynamoAV()
	if err != nil {
		return application.NewRequestError(500, "Failed to marshal UserRole", err)
	}

	// write item to dynamo
	_, err = repo.DDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           repo.TableName,
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(domain_id) AND attribute_not_exists(user_id)"),
	})

	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return application.NewRequestError(400, "Role already exists", nil)
		}

		return application.NewRequestError(500, "Failed to create UserRole", err)
	}

	return nil
}

func (repo *UserRolesRepository) Update(ctx context.Context, userRole *entities.UserRole) *application.RequestError {
	// marshal entity
	userRole.UpdatedAt.Touch()
	item, err := userRole.MarshalDynamoAV()
	if err != nil {
		return application.NewRequestError(500, "Failed to marshal UserRole", err)
	}

	// write item to dynamo
	_, err = repo.DDBClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: repo.TableName,
		Key: map[string]ddbTypes.AttributeValue{
			"domain_id": item["domain_id"],
			"user_id":   item["user_id"],
		},
		UpdateExpression: aws.String("set role_name = :role, updated_at = :ts"),
		ExpressionAttributeValues: map[string]ddbTypes.AttributeValue{
			":role": item["role_name"],
			":ts":   item["updated_at"],
		},
		// only allow updates for existing users
		ConditionExpression: aws.String("attribute_exists(domain_id) AND attribute_exists(user_id)"),
	})

	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return application.NewRequestError(400, "You can only update existing roles", nil)
		}

		return application.NewRequestError(500, "Failed to update UserRole", nil)
	}

	return nil
}
