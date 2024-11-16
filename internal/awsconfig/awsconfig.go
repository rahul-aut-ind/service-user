package awsconfig

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	environ "github.com/rahul-aut-ind/service-user/internal/config"
)

type AWSConfig struct {
	Config *aws.Config
}

// nolint // already aware of deprecated warning
func NewAWSConfig(env *environ.Env) *AWSConfig {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(env.DefaultAWSRegion),
	)
	if env.Environment == environ.LocalEnvironment {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(env.DefaultAWSRegion),
			config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     env.AwsAccessKey,
					SecretAccessKey: env.AwsSecretAccessKey,
				}}),
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           env.DynamoDBConnectionString,
						SigningRegion: env.DefaultAWSRegion,
					}, nil
				})),
		)
	}
	if err != nil {
		log.Fatalf("unable to load AWS config, %v", err)
	}
	return &AWSConfig{Config: &cfg}
}
