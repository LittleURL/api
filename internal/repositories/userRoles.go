package repositories

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	// "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gitlab.com/deltabyte_/littleurl/api/internal/application"
	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
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

func (repo *UserRolesRepository) Insert(ctx context.Context, userRole *entities.UserRole) error {
	// marshal entity
	item, err := userRole.MarshalDynamoAV()
	if err != nil {
		return err
	}

	// write item to dynamo
	_, err = repo.DDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           repo.TableName,
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(domain_id) AND attribute_not_exists(user_id)"),
	})
	return err
}
