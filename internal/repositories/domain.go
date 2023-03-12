package repositories

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gitlab.com/deltabyte_/littleurl/api/internal/application"
	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
)

type DomainsRepository struct {
	TableName *string
	DDBClient *dynamodb.Client
}

func NewDomainsRepository(app *application.Application) *DomainsRepository {
	return &DomainsRepository{
		TableName: aws.String(app.Cfg.Tables.Domains),
		DDBClient: app.DDBClient,
	}
}

func (repo *DomainsRepository) Find(id entities.DomainID) (*entities.Domain, *application.RequestError) {
	res, err := repo.DDBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: repo.TableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
	})
	if err != nil || res == nil || res.Item == nil {
		return nil, application.NewRequestError(404, "Domain not found", err)
	}

	// unmarshal
	domain := &entities.Domain{}
	if err := domain.UnmarshalDynamoAV(res.Item); err != nil {
		return nil, application.NewRequestError(500, "Failed to unmarshal Domain", err)
	}

	return domain, nil
}
