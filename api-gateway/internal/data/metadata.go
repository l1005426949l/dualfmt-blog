package data

import (
	"context"
	"fmt"
)

// mockArticleRepo is a mock implementation of ArticleRepo.
type mockArticleRepo struct{}

// GetArticle returns a single mock article.
func (r *mockArticleRepo) GetArticle(ctx context.Context, id string) (*Article, error) {
	if id == "1" {
		return &Article{
			ID:      "1",
			Title:   "Hello from Data Layer",
			Content: "This content will be fetched by the FileRepo.", // Placeholder
			Version: "1.0.0",
			Format:  "markdown",
			Tags:    []string{"go", "kratos"},
		}, nil
	}
	return nil, fmt.Errorf("article not found: %s", id)
}

// ListArticles returns a list of mock articles.
func (r *mockArticleRepo) ListArticles(ctx context.Context, tag string, limit, offset int) ([]*Article, error) {
	allArticles := []*Article{
		{
			ID:      "1",
			Title:   "Hello from Data Layer",
			Content: "...",
			Version: "1.0.0",
			Format:  "markdown",
			Tags:    []string{"go", "kratos"},
		},
		{
			ID:      "2",
			Title:   "Typst Support in Data Layer",
			Content: "...",
			Version: "1.0.1",
			Format:  "typst",
			Tags:    []string{"typst", "rust"},
		},
	}
	return allArticles, nil
}
