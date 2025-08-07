package data

import (
	"context"
	// Let's assume these generated packages exist for now.
	// I will define the interfaces manually to make the code compile.
	// metadatav1 "api-gateway/api/metadata/v1"
	// filemgmtv1 "api-gateway/api/filemanagement/v1"
)

// ArticleRepo is the repository interface for article metadata.
type ArticleRepo interface {
	GetArticle(ctx context.Context, id string) (*Article, error)
	ListArticles(ctx context.Context, tag string, limit, offset int) ([]*Article, error)
}

// FileRepo is the repository interface for file management.
type FileRepo interface {
	UploadFile(ctx context.Context, filename string, content []byte) (string, error)
	GetFile(ctx context.Context, path string) ([]byte, error)
}

// Article is a simplified domain model for an article.
// It combines data that might come from multiple services.
type Article struct {
	ID      string
	Title   string
	Content string
	Version string
	Format  string
	Tags    []string
}

// Data struct holds all the repository implementations.
// It would normally hold database connections, gRPC clients, etc.
type Data struct {
	// For now, we embed the repositories directly.
	// In a real app, you might have concrete repo structs that hold clients.
	articleRepo ArticleRepo
	fileRepo    FileRepo
}

// NewData creates a new Data object with mock repository implementations.
func NewData() (*Data, func(), error) {
	// This is where you would initialize real gRPC clients.
	// For now, we are creating mock repositories.
	mockArticleRepo := &mockArticleRepo{}
	mockFileRepo := &mockFileRepo{}

	d := &Data{
		articleRepo: mockArticleRepo,
		fileRepo:    mockFileRepo,
	}

	// The cleanup function is used to close connections, etc.
	cleanup := func() {
		// No-op for now
	}

	return d, cleanup, nil
}
