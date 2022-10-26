package application

import (
	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gitlab.com/deltabyte_/littleurl/api/internal/config"
)

type Application struct {
	Ctx       context.Context
	Cfg       *config.Config
	AwsCfg    awsConfig.Config
	DDBClient *dynamodb.Client
}

func New(ctx context.Context) (*Application, error) {
	cfg := config.Load()

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	ddb := dynamodb.NewFromConfig(awsCfg)

	return &Application{Ctx: ctx, Cfg: cfg, AwsCfg: awsCfg, DDBClient: ddb}, nil
}
