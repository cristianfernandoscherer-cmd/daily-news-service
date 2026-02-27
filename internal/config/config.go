package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AWSRegion         string
	DynamoDBTableName string
	DynamoDBEndpoint  string
	BedrockModelID    string
	LogLevel          string
	TopNArticles      int
	UseMockLLM        bool
	Environment       string
}

func Load() (*Config, error) {
	region := os.Getenv("AWS_REGION_NAME")
	if region == "" {
		region = os.Getenv("AWS_REGION")
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = os.Getenv("ENV")
	}

	cfg := &Config{
		AWSRegion:         region,
		DynamoDBTableName: os.Getenv("DYNAMODB_TABLE_NAME"),
		DynamoDBEndpoint:  os.Getenv("DYNAMODB_ENDPOINT"),
		BedrockModelID:    os.Getenv("BEDROCK_MODEL_ID"),
		LogLevel:          os.Getenv("LOG_LEVEL"),
		UseMockLLM:        os.Getenv("USE_MOCK_LLM") == "true",
		Environment:       env,
	}

	topNStr := os.Getenv("TOP_N_ARTICLES")
	if topNStr != "" {
		topN, err := strconv.Atoi(topNStr)
		if err != nil {
			return nil, fmt.Errorf("invalid TOP_N_ARTICLES value %q: %w", topNStr, err)
		}
		cfg.TopNArticles = topN
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	var missing []string

	if c.AWSRegion == "" {
		missing = append(missing, "AWS_REGION/AWS_REGION_NAME")
	}
	if c.DynamoDBTableName == "" {
		missing = append(missing, "DYNAMODB_TABLE_NAME")
	}
	if c.LogLevel == "" {
		missing = append(missing, "LOG_LEVEL")
	}
	if c.Environment == "" {
		missing = append(missing, "ENVIRONMENT/ENV")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}

	if c.TopNArticles <= 0 {
		return fmt.Errorf("TOP_N_ARTICLES must be a positive integer")
	}
	if c.BedrockModelID == "" && !c.UseMockLLM {
		return fmt.Errorf("BEDROCK_MODEL_ID is required when not using mock LLM")
	}
	return nil
}
