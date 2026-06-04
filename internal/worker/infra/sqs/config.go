package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/clodoaldomarques/balances-api/config"
	"github.com/clodoaldomarques/core-sdk/pkg/aws"
	"github.com/clodoaldomarques/core-sdk/pkg/logger"
)

func NewSQSClient(ctx context.Context) *sqs.Client {
	c := config.New()
	cfg, err := aws.NewCustomConfig(ctx, c.AwsRegion, c.AwsAddress, c.AccessKeyID, c.SecretAccessKey)
	if err != nil {
		logger.Fatal(ctx, "falha ao carregar configuração:", logger.Fields{
			"error":      err.Error(),
			"AwsRegion":  cfg.Region,
			"AwsAddress": cfg.BaseEndpoint,
		})
	}
	return sqs.NewFromConfig(cfg)
}
