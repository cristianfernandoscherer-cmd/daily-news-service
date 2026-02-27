package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"time"

	"github.com/cristian-scherer/daily-news-service/internal/news/models"
	"github.com/mmcdole/gofeed"
)

type feedConfig struct {
	Source models.Source
	URL    string
}

var defaultFeeds = []feedConfig{
	{Source: models.SourceDevTo, URL: "https://dev.to/feed"},
	{Source: models.SourceAWSBlog, URL: "https://aws.amazon.com/blogs/aws/feed/"},
	{Source: models.SourceTechCrunch, URL: "https://techcrunch.com/feed/"},
	{Source: models.SourceTheVerge, URL: "https://www.theverge.com/rss/index.xml"},
	{Source: models.SourceWired, URL: "https://www.wired.com/feed/rss"},
	{Source: models.SourceArsTechnica, URL: "https://arstechnica.com/feed/"},
	{Source: models.SourceCNET, URL: "https://www.cnet.com/rss/news/"},
	{Source: models.SourceZDNet, URL: "https://www.zdnet.com/news/rss.xml"},
	{Source: models.SourceVentureBeat, URL: "https://venturebeat.com/feed/"},
	{Source: models.SourceEngadget, URL: "https://www.engadget.com/rss.xml"},
	{Source: models.SourceMITReview, URL: "https://www.technologyreview.com/feed/"},
	{Source: models.SourceHackerNews, URL: "https://news.ycombinator.com/rss"},
}

type FetcherService struct {
	parser *gofeed.Parser
	feeds  []feedConfig
	logger *slog.Logger
}

func NewFetcherService(logger *slog.Logger) *FetcherService {
	return &FetcherService{
		parser: gofeed.NewParser(),
		feeds:  defaultFeeds,
		logger: logger,
	}
}

func (f *FetcherService) Fetch(ctx context.Context) ([]*models.Article, error) {
	seen := make(map[string]struct{})
	var articles []*models.Article

	for _, feed := range f.feeds {
		fetched, err := f.fetchFeed(ctx, feed)
		if err != nil {
			f.logger.Warn("failed to fetch feed, skipping",
				slog.String("source", string(feed.Source)),
				slog.String("url", feed.URL),
				slog.String("error", err.Error()),
			)
			continue
		}

		for _, a := range fetched {
			if _, ok := seen[a.ID]; ok {
				continue
			}
			seen[a.ID] = struct{}{}
			articles = append(articles, a)
		}

		f.logger.Info("fetched articles from source",
			slog.String("source", string(feed.Source)),
			slog.Int("count", len(fetched)),
		)
	}

	f.logger.Info("total articles fetched", slog.Int("total", len(articles)))
	return articles, nil
}

func (f *FetcherService) fetchFeed(ctx context.Context, cfg feedConfig) ([]*models.Article, error) {
	feed, err := f.parser.ParseURLWithContext(cfg.URL, ctx)
	if err != nil {
		return nil, fmt.Errorf("parsing feed %q: %w", cfg.URL, err)
	}

	now := time.Now().UTC()
	var articles []*models.Article

	for _, item := range feed.Items {
		if item.Link == "" {
			continue
		}

		if item.PublishedParsed == nil {
			continue
		}

		pubDate := item.PublishedParsed.UTC()
		if pubDate.Year() != now.Year() || pubDate.Month() != now.Month() || pubDate.Day() != now.Day() {
			continue
		}

		id := generateID(item.Link)
		a := &models.Article{
			ID:          id,
			Title:       item.Title,
			URL:         item.Link,
			Source:      cfg.Source,
			RawContent:  item.Description,
			FetchedAt:   now,
			PublishedAt: &pubDate,
			TTL:         now.Add(30 * 24 * time.Hour).Unix(),
		}

		articles = append(articles, a)
	}

	return articles, nil
}

func generateID(url string) string {
	h := sha256.Sum256([]byte(url))
	return fmt.Sprintf("%x", h[:8])
}
