package helpers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gitlab.com/deltabyte_/littleurl/api/internal/application"
	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
	"gitlab.com/deltabyte_/littleurl/api/internal/permissions"
)

func FindDomainUser(app *application.Application, domainId string, userId string) (*entities.DomainUser, error) {
	res, err := app.DDBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &app.Cfg.Tables.DomainUsers,
		Key: map[string]types.AttributeValue{
			"domain_id": &types.AttributeValueMemberS{
				Value: domainId,
			},
			"user_id": &types.AttributeValueMemberS{
				Value: userId,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// unmarshal domain user, default permission to nobody if the record was not found
	domainUser := &entities.DomainUser{Permission: permissions.Nobody}
	if res != nil && res.Item != nil {
		domainUser.UnmarshalDynamoAV(res.Item)
	}

	return domainUser, nil
}
