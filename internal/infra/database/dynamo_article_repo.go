package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/cristian-scherer/daily-news-service/internal/news/models"
)

type DatabaseClient interface {
	PutItem(ctx context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Query(ctx context.Context, input *dynamodb.QueryInput, opts ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	Scan(ctx context.Context, input *dynamodb.ScanInput, opts ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

type DynamoArticleRepository struct {
	db        DatabaseClient
	tableName string
	logger    *slog.Logger
}

func NewDynamoArticleRepository(db DatabaseClient, tableName string, logger *slog.Logger) *DynamoArticleRepository {
	return &DynamoArticleRepository{db: db, tableName: tableName, logger: logger}
}

func (r *DynamoArticleRepository) Save(ctx context.Context, article *models.Article) error {
	item, err := attributevalue.MarshalMap(article)
	if err != nil {
		return fmt.Errorf("marshalling article %q: %w", article.ID, err)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("saving article %q: %w", article.ID, err)
	}

	r.logger.Debug("article saved", slog.String("id", article.ID), slog.String("title", article.Title))
	return nil
}

func (r *DynamoArticleRepository) SaveBatch(ctx context.Context, articles []*models.Article) error {
	var errs []error

	for _, a := range articles {
		if err := r.Save(ctx, a); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("batch save had %d error(s); first error: %w", len(errs), errs[0])
	}

	r.logger.Info("batch save complete", slog.Int("count", len(articles)))
	return nil
}

func (r *DynamoArticleRepository) GetLatest(ctx context.Context, limit int) ([]*models.Article, error) {
	out, err := r.db.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
		Limit:     aws.Int32(int32(limit * 2)),
	})
	if err != nil {
		return nil, fmt.Errorf("scanning articles table: %w", err)
	}

	var articles []*models.Article
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &articles); err != nil {
		return nil, fmt.Errorf("unmarshalling articles: %w", err)
	}

	cutoff := time.Now().UTC().Add(-24 * time.Hour)
	var recent []*models.Article
	for _, a := range articles {
		if a.FetchedAt.After(cutoff) {
			recent = append(recent, a)
		}
	}

	if len(recent) > limit {
		recent = recent[:limit]
	}

	return recent, nil
}
