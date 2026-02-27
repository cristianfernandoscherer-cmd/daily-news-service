package di

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/cristian-scherer/daily-news-service/internal/config"
	infra "github.com/cristian-scherer/daily-news-service/internal/infra/database"
	"github.com/cristian-scherer/daily-news-service/internal/news/handler"
	"github.com/cristian-scherer/daily-news-service/internal/news/service"
	pkgdb "github.com/cristian-scherer/daily-news-service/pkg/database"
	"github.com/cristian-scherer/daily-news-service/pkg/llm"
)

type Container struct {
	Config            *config.Config
	UpdateNewsHandler *handler.UpdateNewsHandler
}

func Build(ctx context.Context) (*Container, error) {

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	db, err := pkgdb.NewDynamoDBClient(ctx, cfg.AWSRegion, cfg.DynamoDBEndpoint)
	if err != nil {
		return nil, fmt.Errorf("creating DynamoDB client: %w", err)
	}

	var llmClient llm.LLMClient
	if cfg.Environment == "prod" {
		var err error
		llmClient, err = llm.NewBedrockClient(ctx, cfg.AWSRegion, cfg.BedrockModelID)
		if err != nil {
			return nil, fmt.Errorf("creating Bedrock client: %w", err)
		}
	} else {
		llmClient = service.NewMockLLMClient(logger)
	}

	articleRepo := infra.NewDynamoArticleRepository(db, cfg.DynamoDBTableName, logger)

	fetcherSvc := service.NewFetcherService(logger)
	summarizerSvc := service.NewSummarizerService(llmClient, logger)
	rankerSvc := service.NewRankerService(cfg.TopNArticles, logger)

	updateNewsHandler := handler.NewUpdateNewsHandler(
		fetcherSvc,
		summarizerSvc,
		rankerSvc,
		articleRepo,
		logger,
	)

	return &Container{
		Config:            cfg,
		UpdateNewsHandler: updateNewsHandler,
	}, nil
}
