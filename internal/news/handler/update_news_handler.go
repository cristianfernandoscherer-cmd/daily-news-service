package handler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cristian-scherer/daily-news-service/internal/news/repository"
	"github.com/cristian-scherer/daily-news-service/internal/news/service"
)

type UpdateNewsHandler struct {
	fetcher    *service.FetcherService
	summarizer *service.SummarizerService
	ranker     *service.RankerService
	repo       repository.ArticleRepository
	logger     *slog.Logger
}

func NewUpdateNewsHandler(
	fetcher *service.FetcherService,
	summarizer *service.SummarizerService,
	ranker *service.RankerService,
	repo repository.ArticleRepository,
	logger *slog.Logger,
) *UpdateNewsHandler {
	return &UpdateNewsHandler{
		fetcher:    fetcher,
		summarizer: summarizer,
		ranker:     ranker,
		repo:       repo,
		logger:     logger,
	}
}

func (h *UpdateNewsHandler) Handle(ctx context.Context) error {
	h.logger.Info("starting daily news update.")

	articles, err := h.fetcher.Fetch(ctx)
	if err != nil {
		return fmt.Errorf("fetching articles: %w", err)
	}
	if len(articles) == 0 {
		h.logger.Warn("no articles fetched; skipping run")
		return nil
	}
	h.logger.Info("fetch complete", slog.Int("articles", len(articles)))

	resumed, err := h.summarizer.ResumeArticles(ctx, articles)
	if err != nil {
		return fmt.Errorf("resuming articles: %w", err)
	}
	if len(resumed) == 0 {
		h.logger.Warn("no articles successfully resumed; skipping run")
		return nil
	}
	h.logger.Info("resumption complete", slog.Int("resumed", len(resumed)))

	top := h.ranker.Rank(resumed)
	h.logger.Info("ranking complete", slog.Int("selected", len(top)))

	if err := h.repo.SaveBatch(ctx, top); err != nil {
		return fmt.Errorf("saving articles to database: %w", err)
	}

	h.logger.Info("daily news update complete", slog.Int("saved", len(top)))

	return nil
}
