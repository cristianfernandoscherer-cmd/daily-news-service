package repository

import (
	"context"

	"github.com/cristian-scherer/daily-news-service/internal/news/models"
)

type ArticleRepository interface {
	Save(ctx context.Context, article *models.Article) error
	SaveBatch(ctx context.Context, articles []*models.Article) error
	GetLatest(ctx context.Context, limit int) ([]*models.Article, error)
}
