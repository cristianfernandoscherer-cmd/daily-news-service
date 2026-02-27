package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type BedrockClient struct {
	client  *bedrockruntime.Client
	modelID string
}

func NewBedrockClient(ctx context.Context, region, modelID string) (*BedrockClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("loading AWS config for Bedrock: %w", err)
	}

	client := bedrockruntime.NewFromConfig(cfg)
	return &BedrockClient{client: client, modelID: modelID}, nil
}

type claudeRequest struct {
	AnthropicVersion string          `json:"anthropic_version"`
	MaxTokens        int             `json:"max_tokens"`
	Messages         []claudeMessage `json:"messages"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

func (b *BedrockClient) ResumeArticle(ctx context.Context, prompt string) (string, error) {
	reqBody := claudeRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        1024,
		Messages: []claudeMessage{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshalling Bedrock request: %w", err)
	}

	resp, err := b.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(b.modelID),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return "", fmt.Errorf("invoking Bedrock model %q: %w", b.modelID, err)
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(resp.Body, &claudeResp); err != nil {
		return "", fmt.Errorf("parsing Bedrock response: %w", err)
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("empty response from Bedrock model")
	}

	return claudeResp.Content[0].Text, nil
}
