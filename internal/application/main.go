package application

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	lumigo "github.com/lumigo-io/lumigo-go-tracer"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gitlab.com/deltabyte_/littleurl/api/internal/config"
)

type Application struct {
	Ctx       context.Context
	Cfg       *config.Config
	AwsCfg    aws.Config
	DDBClient *dynamodb.Client
}

func New(ctx context.Context) (*Application, error) {
	cfg := config.Load()

	// lumigo
	tracedClient := &http.Client{
    Transport: lumigo.NewTransport(http.DefaultTransport),
  }

	// aws config
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithHTTPClient(tracedClient))
	if err != nil {
		return nil, err
	}

	ddb := dynamodb.NewFromConfig(awsCfg)

	return &Application{Ctx: ctx, Cfg: cfg, AwsCfg: awsCfg, DDBClient: ddb}, nil
}
