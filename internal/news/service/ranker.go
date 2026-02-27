package service

import (
	"log/slog"
	"sort"

	"github.com/cristian-scherer/daily-news-service/internal/news/models"
)

type RankerService struct {
	topN   int
	logger *slog.Logger
}

func NewRankerService(topN int, logger *slog.Logger) *RankerService {
	return &RankerService{topN: topN, logger: logger}
}

func (r *RankerService) Rank(articles []*models.Article) []*models.Article {
	ranked := make([]*models.Article, len(articles))
	copy(ranked, articles)

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].RelevanceScore > ranked[j].RelevanceScore
	})

	if len(ranked) <= r.topN {
		r.logger.Info("ranking complete",
			slog.Int("total", len(ranked)),
			slog.Int("selected", len(ranked)),
		)
		return ranked
	}

	selected := ranked[:r.topN]
	r.logger.Info("ranking complete",
		slog.Int("total", len(articles)),
		slog.Int("selected", len(selected)),
		slog.Float64("min_score", selected[len(selected)-1].RelevanceScore),
	)

	return selected
}
