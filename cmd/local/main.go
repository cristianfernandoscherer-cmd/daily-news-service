// cmd/local/main.go
// Local development runner — loads .env and triggers the full news pipeline
// without deploying to AWS. Requires LocalStack running (make infra-up + make migrate).
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/cristian-scherer/daily-news-service/internal/di"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file for local configuration.
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️  Note: .env file not found, using system environment variables")
	}

	// Specifically for local runner, default to 'dev' if not set
	if os.Getenv("ENVIRONMENT") == "" && os.Getenv("ENV") == "" {
		log.Println("ℹ️  ENVIRONMENT not set, defaulting to 'dev' for local run")
		os.Setenv("ENVIRONMENT", "dev")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log.Println("building dependency container...")
	container, err := di.Build(ctx)
	if err != nil {
		log.Fatalf("failed to build container: %v", err)
	}

	log.Println("running news update pipeline...")
	if err := container.UpdateNewsHandler.Handle(ctx); err != nil {
		log.Fatalf("pipeline failed: %v", err)
	}

	log.Println("✅ local run complete")
}
