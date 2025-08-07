package biz

import (
	"context"
	"api-gateway/internal/data"
)

// ArticleUsecase is a usecase for handling article-related logic.
// It orchestrates calls to multiple repositories.
type ArticleUsecase struct {
	articleRepo data.ArticleRepo
	fileRepo    data.FileRepo
}

// NewArticleUsecase creates a new ArticleUsecase.
func NewArticleUsecase(articleRepo data.ArticleRepo, fileRepo data.FileRepo) *ArticleUsecase {
	return &ArticleUsecase{
		articleRepo: articleRepo,
		fileRepo:    fileRepo,
	}
}

// GetArticle orchestrates getting metadata and content for an article.
func (uc *ArticleUsecase) GetArticle(ctx context.Context, id string) (*data.Article, error) {
	// 1. Get metadata from the article repo
	article, err := uc.articleRepo.GetArticle(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. Get content from the file repo
	// In a real system, the article metadata would contain the file path.
	// We'll construct a mock path for now.
	mockPath := "mock/storage/path/for/article-1.md"
	content, err := uc.fileRepo.GetFile(ctx, mockPath)
	if err != nil {
		// We could decide to return the article with empty content or return an error.
		// Let's return the error.
		return nil, err
	}

	// 3. Combine and return
	article.Content = string(content)
	return article, nil
}

// ListArticles simply passes through to the repo for now.
func (uc *ArticleUsecase) ListArticles(ctx context.Context, tag string, limit, offset int) ([]*data.Article, error) {
	// In a real app, you might enrich the articles with more data here.
	return uc.articleRepo.ListArticles(ctx, tag, limit, offset)
}
