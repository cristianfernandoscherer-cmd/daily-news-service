package models

import "time"

type Source string

const (
	SourceDevTo       Source = "devto"
	SourceAWSBlog     Source = "aws_blog"
	SourceTechCrunch  Source = "techcrunch"
	SourceTheVerge    Source = "the_verge"
	SourceWired       Source = "wired"
	SourceArsTechnica Source = "ars_technica"
	SourceCNET        Source = "cnet"
	SourceZDNet       Source = "zdnet"
	SourceVentureBeat Source = "venturebeat"
	SourceEngadget    Source = "engadget"
	SourceMITReview   Source = "mit_review"
	SourceHackerNews  Source = "hacker_news"
)

type Article struct {
	ID             string     `json:"id" dynamodbav:"id"`
	Title          string     `json:"title" dynamodbav:"title"`
	URL            string     `json:"url" dynamodbav:"url"`
	Source         Source     `json:"source" dynamodbav:"source"`
	RawContent     string     `json:"raw_content,omitempty" dynamodbav:"-"`
	Summary        string     `json:"summary" dynamodbav:"summary"`
	KeyPoints      []string   `json:"key_points" dynamodbav:"key_points"`
	RelevanceScore float64    `json:"relevance_score" dynamodbav:"relevance_score"`
	FetchedAt      time.Time  `json:"fetched_at" dynamodbav:"fetched_at"`
	PublishedAt    *time.Time `json:"published_at,omitempty" dynamodbav:"published_at,omitempty"`
	TTL            int64      `json:"ttl" dynamodbav:"ttl"`
}
