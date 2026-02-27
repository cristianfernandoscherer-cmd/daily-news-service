// .scripts/migrate/migrate.go
// Creates the DynamoDB articles table in LocalStack for local development.
// Run via: go run ./.scripts/migrate/migrate.go
// Or:      make migrate
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pkgdb "github.com/cristian-scherer/daily-news-service/pkg/database"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	region := getEnv("AWS_REGION_NAME", "us-east-1")
	endpoint := getEnv("DYNAMODB_ENDPOINT", "http://localhost:4566")
	tableName := getEnv("DYNAMODB_TABLE_NAME", "articles")

	ctx := context.Background()

	log.Printf("→ connecting to DynamoDB at %s (region: %s)", endpoint, region)
	db, err := pkgdb.NewDynamoDBClient(ctx, region, endpoint)
	if err != nil {
		log.Fatalf("failed to create DynamoDB client: %v", err)
	}

	log.Printf("→ creating table: %s", tableName)
	if err := db.CreateTable(ctx, tableName); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Printf("✅ table %q created successfully\n", tableName)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
