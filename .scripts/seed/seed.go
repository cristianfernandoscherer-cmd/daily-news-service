// .scripts/seed/seed.go
// Seeds sample article data into LocalStack DynamoDB for local development.
// Run via: go run ./.scripts/seed/seed.go
// Or:      make seed
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cristian-scherer/daily-news-service/internal/infra/database"
	"github.com/cristian-scherer/daily-news-service/internal/news/models"
	pkgdb "github.com/cristian-scherer/daily-news-service/pkg/database"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	region := getEnv("AWS_REGION_NAME", "us-east-1")
	endpoint := getEnv("DYNAMODB_ENDPOINT", "http://localhost:4566")
	tableName := getEnv("DYNAMODB_TABLE_NAME", "articles")

	ctx := context.Background()

	db, err := pkgdb.NewDynamoDBClient(ctx, region, endpoint)
	if err != nil {
		log.Fatalf("failed to create DynamoDB client: %v", err)
	}

	repo := database.NewDynamoArticleRepository(db, tableName, nil)

	now := time.Now().UTC()
	articles := []*models.Article{
		{
			ID:             "seed001",
			Title:          "Go 1.23 Released with New Iterators",
			URL:            "https://go.dev/blog/go1.23",
			Source:         models.SourceDevTo,
			Summary:        "Go 1.23 introduces first-class iterator support via range-over-func, simplifying collection traversal.",
			KeyPoints:      []string{"range-over-func support", "improved standard library", "minor performance gains"},
			RelevanceScore: 9.2,
			FetchedAt:      now,
			TTL:            now.Add(30 * 24 * time.Hour).Unix(),
		},
		{
			ID:             "seed002",
			Title:          "AWS Introduces Bedrock RAG Improvements",
			URL:            "https://aws.amazon.com/blogs/aws/bedrock-rag-2024",
			Source:         models.SourceAWSBlog,
			Summary:        "Amazon Bedrock now supports improved retrieval-augmented generation pipelines with native knowledge base connectors.",
			KeyPoints:      []string{"native RAG support", "knowledge base integration", "reduced latency"},
			RelevanceScore: 8.7,
			FetchedAt:      now,
			TTL:            now.Add(30 * 24 * time.Hour).Unix(),
		},
		{
			ID:             "seed003",
			Title:          "GitHub Copilot Gets Multi-File Editing",
			URL:            "https://github.blog/copilot-multi-file",
			Source:         models.SourceGitHub,
			Summary:        "GitHub Copilot now supports editing multiple files simultaneously, enabling broader refactoring tasks directly in the IDE.",
			KeyPoints:      []string{"multi-file editing", "refactoring support", "IDE integration"},
			RelevanceScore: 8.1,
			FetchedAt:      now,
			TTL:            now.Add(30 * 24 * time.Hour).Unix(),
		},
	}

	log.Printf("→ seeding %d articles into %s", len(articles), tableName)
	if err := repo.SaveBatch(ctx, articles); err != nil {
		log.Fatalf("seeding failed: %v", err)
	}

	fmt.Printf("✅ seeded %d articles successfully\n", len(articles))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
