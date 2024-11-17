package integrationtest

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/testcontainers/testcontainers-go"
	tcdynamodb "github.com/testcontainers/testcontainers-go/modules/dynamodb"
)

const (
	UserImageTable = "user-images"
	HashKey        = "UserID"
	RangeKey       = "ImageID"
	IndexRangeKey  = "TakenAt"
	SecondaryIndex = "UserIDTakenAtIndex"
)

type (
	DynamoDBSetup struct {
		container *tcdynamodb.DynamoDBContainer
		Client    *dynamodb.Client
	}

	dynamoDBResolver struct {
		HostPort string
	}
)

func NewDynamoDbSetup() *DynamoDBSetup {
	setup := &DynamoDBSetup{}
	setup.initialize()
	return setup
}

func (tc *DynamoDBSetup) initialize() {
	ctx := context.Background()

	err := tc.createDynamoDBContainer(ctx)
	if err != nil {
		panic(err)
	}
	err = tc.initDynamoDBClient(ctx)
	if err != nil {
		panic(err)
	}
	err = tc.createTable()
	if err != nil {
		panic(err)
	}
}

func (tc *DynamoDBSetup) createDynamoDBContainer(ctx context.Context) error {
	log.Printf("starting test container")
	ctr, err := tcdynamodb.Run(ctx, "amazon/dynamodb-local:2.2.1")
	if err != nil {
		return fmt.Errorf("failed to run dynamodb container: %s", err)
	}
	tc.container = ctr
	log.Printf("test container started")
	return nil
}

func (tc *DynamoDBSetup) initDynamoDBClient(ctx context.Context) error {
	hostPort, err := tc.container.ConnectionString(ctx)
	if err != nil {
		return fmt.Errorf("failed to get dynamodb conn string: %s", err)
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID:     "dummy",
			SecretAccessKey: "dummy",
		},
	}))
	if err != nil {
		return fmt.Errorf("failed to create dynamodb client: %s", err)
	}

	tc.Client = dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolverV2(&dynamoDBResolver{HostPort: hostPort}))
	return nil
}

func (tc *DynamoDBSetup) createTable() error {
	_, err := tc.Client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: aws.String(UserImageTable),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(HashKey),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String(RangeKey),
				KeyType:       types.KeyTypeRange,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(HashKey),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(RangeKey),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(IndexRangeKey),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		LocalSecondaryIndexes: []types.LocalSecondaryIndex{
			{
				IndexName: aws.String(SecondaryIndex),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String(HashKey), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String(IndexRangeKey), KeyType: types.KeyTypeRange},
				},
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func (r *dynamoDBResolver) ResolveEndpoint(ctx context.Context, params dynamodb.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return smithyendpoints.Endpoint{
		URI: url.URL{Host: r.HostPort, Scheme: "http"},
	}, nil
}

func (tc *DynamoDBSetup) Stop() {
	if err := testcontainers.TerminateContainer(tc.container); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}
