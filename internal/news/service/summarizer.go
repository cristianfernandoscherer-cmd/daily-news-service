package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/cristian-scherer/daily-news-service/internal/news/models"
	pkgllm "github.com/cristian-scherer/daily-news-service/pkg/llm"
)

type llmResponse struct {
	Summary        string   `json:"summary"`
	KeyPoints      []string `json:"key_points"`
	RelevanceScore float64  `json:"relevance_score"`
}

type SummarizerService struct {
	llm    pkgllm.LLMClient
	logger *slog.Logger
}

func NewSummarizerService(llmClient pkgllm.LLMClient, logger *slog.Logger) *SummarizerService {
	return &SummarizerService{llm: llmClient, logger: logger}
}

func (s *SummarizerService) ResumeArticles(ctx context.Context, articles []*models.Article) ([]*models.Article, error) {
	var resumed []*models.Article

	for _, a := range articles {
		if err := ctx.Err(); err != nil {
			return resumed, err
		}

		result, err := s.processSingleArticle(ctx, a)
		if err != nil {
			s.logger.Warn("failed to resume article, skipping",
				slog.String("article_id", a.ID),
				slog.String("title", a.Title),
				slog.String("error", err.Error()),
			)
			continue
		}

		a.Summary = result.Summary
		a.KeyPoints = result.KeyPoints
		a.RelevanceScore = result.RelevanceScore
		resumed = append(resumed, a)
	}

	return resumed, nil
}

func (s *SummarizerService) processSingleArticle(ctx context.Context, a *models.Article) (*llmResponse, error) {
	prompt := getPrompt(a)

	rawResponse, err := s.llm.ResumeArticle(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("llm call: %w", err)
	}

	var resp llmResponse
	if err := json.Unmarshal([]byte(rawResponse), &resp); err != nil {
		return nil, fmt.Errorf("parsing llm response JSON: %w (raw: %q)", err, rawResponse)
	}

	if resp.RelevanceScore < 0 || resp.RelevanceScore > 10 {
		return nil, fmt.Errorf("relevance_score %v out of expected 0–10 range", resp.RelevanceScore)
	}

	return &resp, nil
}

func getPrompt(a *models.Article) string {
	return fmt.Sprintf(`Você é um engenheiro de software sênior e curador de notícias de tecnologia.
Analise o seguinte artigo e responda APENAS com um objeto JSON válido — sem markdown, sem preâmbulo.

Título do Artigo: %s
URL do Artigo: %s
Conteúdo / Descrição do Artigo:
%s

Retorne exatamente este formato JSON:
{
  "summary": "<resumo de no máximo 5 linhas em Português do Brasil>",
  "key_points": ["<ponto 1>", "<ponto 2>", "<ponto 3>"],
  "relevance_score": <float entre 0 e 10, onde 10 é extremamente relevante para engenheiros de software>
}`, a.Title, a.URL, a.RawContent)
}
