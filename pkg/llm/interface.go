package llm

import "context"

type LLMClient interface {
	ResumeArticle(ctx context.Context, prompt string) (string, error)
}
