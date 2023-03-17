package repositories

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gitlab.com/deltabyte_/littleurl/api/internal/application"
	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
)

const userInvitesIndexDomain = "domain-invites"

type UserInvitesRepository struct {
	TableName *string
	DDBClient *dynamodb.Client
}

func NewUserInvitesRepository(app *application.Application) *UserInvitesRepository {
	return &UserInvitesRepository{
		TableName: aws.String(app.Cfg.Tables.UserInvites),
		DDBClient: app.DDBClient,
	}
}

func (repo *UserInvitesRepository) Find(ctx context.Context, id entities.UserInviteID) (*entities.UserInvite, *application.RequestError) {
	res, err := repo.DDBClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: repo.TableName,
		Key: map[string]ddbTypes.AttributeValue{
			"id": &ddbTypes.AttributeValueMemberS{
				Value: id,
			},
		},
	})
	if err != nil || res == nil || res.Item == nil {
		return nil, application.NewRequestError(404, "UserInvite not found", err)
	}

	// unmarshal
	domain := &entities.UserInvite{}
	if err := domain.UnmarshalDynamoAV(res.Item); err != nil {
		return nil, application.NewRequestError(500, "Failed to unmarshal UserInvite", err)
	}

	return domain, nil
}

func (repo *UserInvitesRepository) Cleanup(ctx context.Context, userInvite *entities.UserInvite) *application.RequestError {
	// find all invites for a domain, filtered by email
	res, err := repo.DDBClient.Query(ctx, &dynamodb.QueryInput{
		TableName: repo.TableName,
		IndexName: aws.String(userInvitesIndexDomain),
		KeyConditionExpression: aws.String("domain_id = :id"),
		FilterExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]ddbTypes.AttributeValue{
			":id": &ddbTypes.AttributeValueMemberS{
				Value: userInvite.DomainID,
			},
			":email": &ddbTypes.AttributeValueMemberS{
				Value: userInvite.Email,
			},
		},
	})
	if err != nil || res == nil || len(res.Items) == 0 {
		return application.NewRequestError(500, "Failed to load UserInvites by Domain", err)
	}

	// unmarshal
	invites := entities.UserInvites{}
	if err := invites.UnmarshalDynamoAV(res.Items); err != nil {
		return application.NewRequestError(500, "Failed to unmarshal UserInvites", err)
	}

	// create delete requests
	writeRequests := []ddbTypes.WriteRequest{}
	for _, invite := range invites {
		writeRequests = append(writeRequests, ddbTypes.WriteRequest{
			DeleteRequest: &ddbTypes.DeleteRequest{
				Key: map[string]ddbTypes.AttributeValue{
					":id": &ddbTypes.AttributeValueMemberS{
						Value: invite.ID,
					},
				},
			},
		})
	}

	// delete all of the invites
	_, err = repo.DDBClient.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]ddbTypes.WriteRequest{
			*repo.TableName: writeRequests,
		},
	})
	if err != nil {
		return application.NewRequestError(500, "Failed to delete UserInvites", err)
	}

	return nil
}
