package llm

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
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

func (b *BedrockClient) ResumeArticle(ctx context.Context, prompt string) (string, error) {
	input := &bedrockruntime.ConverseInput{
		ModelId: aws.String(b.modelID),
		Messages: []types.Message{
			{
				Role: types.ConversationRoleUser,
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{Value: prompt},
				},
			},
		},
		InferenceConfig: &types.InferenceConfiguration{
			MaxTokens: aws.Int32(2048),
		},
	}

	resp, err := b.client.Converse(ctx, input)
	if err != nil {
		return "", fmt.Errorf("invoking Bedrock model %q via Converse: %w", b.modelID, err)
	}

	output, ok := resp.Output.(*types.ConverseOutputMemberMessage)
	if !ok {
		return "", fmt.Errorf("unexpected output type from Bedrock: %T", resp.Output)
	}

	if len(output.Value.Content) == 0 {
		return "", fmt.Errorf("empty response from Bedrock model")
	}

	textOutput, ok := output.Value.Content[0].(*types.ContentBlockMemberText)
	if !ok {
		return "", fmt.Errorf("unexpected content block type: %T", output.Value.Content[0])
	}

	return textOutput.Value, nil
}
