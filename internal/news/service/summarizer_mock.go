package service

import (
	"context"
	"encoding/json"
	"log/slog"
)

type MockLLMClient struct {
	logger *slog.Logger
}

func NewMockLLMClient(logger *slog.Logger) *MockLLMClient {
	return &MockLLMClient{logger: logger}
}

func (m *MockLLMClient) ResumeArticle(ctx context.Context, prompt string) (string, error) {
	m.logger.Debug("mock llm called", slog.String("prompt_preview", prompt[:50]+"..."))

	resp := llmResponse{
		Summary:        "Este é um resumo gerado pelo mock para testes locais. O sistema agora suporta múltiplos feeds e filtra por data corretamente.",
		KeyPoints:      []string{"Ponto 1: Sem necessidade de AWS Bedrock", "Ponto 2: Fluxo puramente local", "Ponto 3: Execução rápida"},
		RelevanceScore: 8.5,
	}

	raw, _ := json.Marshal(resp)
	return string(raw), nil
}
