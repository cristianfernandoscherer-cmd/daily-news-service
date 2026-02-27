package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cristian-scherer/daily-news-service/internal/di"
)

var container *di.Container

func init() {
	var err error
	container, err = di.Build(context.Background())
	if err != nil {
		log.Fatalf("failed to initialise container: %v", err)
	}
}

func handler(ctx context.Context) error {
	return container.UpdateNewsHandler.Handle(ctx)
}

func main() {
	lambda.Start(handler)
}
