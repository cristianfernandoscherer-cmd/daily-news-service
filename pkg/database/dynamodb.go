package database

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient struct {
	client *dynamodb.Client
}

func NewDynamoDBClient(ctx context.Context, region, endpoint string) (*DynamoDBClient, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	dynamoOpts := []func(*dynamodb.Options){}
	if endpoint != "" {
		dynamoOpts = append(dynamoOpts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}

	client := dynamodb.NewFromConfig(cfg, dynamoOpts...)
	return &DynamoDBClient{client: client}, nil
}

func (d *DynamoDBClient) PutItem(ctx context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	out, err := d.client.PutItem(ctx, input, opts...)
	if err != nil {
		return nil, fmt.Errorf("dynamodb PutItem: %w", err)
	}
	return out, nil
}

func (d *DynamoDBClient) GetItem(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	out, err := d.client.GetItem(ctx, input, opts...)
	if err != nil {
		return nil, fmt.Errorf("dynamodb GetItem: %w", err)
	}
	return out, nil
}

func (d *DynamoDBClient) Query(ctx context.Context, input *dynamodb.QueryInput, opts ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	out, err := d.client.Query(ctx, input, opts...)
	if err != nil {
		return nil, fmt.Errorf("dynamodb Query: %w", err)
	}
	return out, nil
}

func (d *DynamoDBClient) Scan(ctx context.Context, input *dynamodb.ScanInput, opts ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	out, err := d.client.Scan(ctx, input, opts...)
	if err != nil {
		return nil, fmt.Errorf("dynamodb Scan: %w", err)
	}
	return out, nil
}

func (d *DynamoDBClient) CreateTable(ctx context.Context, tableName string) error {
	_, err := d.client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName:   aws.String(tableName),
		BillingMode: types.BillingModePayPerRequest,
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("fetched_at"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			{AttributeName: aws.String("fetched_at"), KeyType: types.KeyTypeRange},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("fetched_at-index"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("fetched_at"), KeyType: types.KeyTypeHash},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("creating table %q: %w", tableName, err)
	}
	return nil
}
