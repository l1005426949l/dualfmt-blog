package biz

import (
	"context"
	"errors"
	"testing"
	"api-gateway/internal/data"
	"github.com/stretchr/testify/assert"
)

// --- Mock Implementations for Testing ---

// mockArticleRepo is a test-specific mock for ArticleRepo.
type mockArticleRepo struct {
	GetArticleFunc   func(ctx context.Context, id string) (*data.Article, error)
	ListArticlesFunc func(ctx context.Context, tag string, limit, offset int) ([]*data.Article, error)
}

func (m *mockArticleRepo) GetArticle(ctx context.Context, id string) (*data.Article, error) {
	return m.GetArticleFunc(ctx, id)
}

func (m *mockArticleRepo) ListArticles(ctx context.Context, tag string, limit, offset int) ([]*data.Article, error) {
	return m.ListArticlesFunc(ctx, tag, limit, offset)
}

// mockFileRepo is a test-specific mock for FileRepo.
type mockFileRepo struct {
	GetFileFunc func(ctx context.Context, path string) ([]byte, error)
	UploadFileFunc func(ctx context.Context, filename string, content []byte) (string, error)
}

func (m *mockFileRepo) GetFile(ctx context.Context, path string) ([]byte, error) {
	return m.GetFileFunc(ctx, path)
}

func (m *mockFileRepo) UploadFile(ctx context.Context, filename string, content []byte) (string, error) {
	return m.UploadFileFunc(ctx, filename, content)
}


// --- Tests for ArticleUsecase ---

func TestArticleUsecase_GetArticle(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockArticle := &data.Article{ID: "1", Title: "Test Title"}
		mockContent := []byte("Test Content")

		articleRepo := &mockArticleRepo{
			GetArticleFunc: func(ctx context.Context, id string) (*data.Article, error) {
				assert.Equal(t, "1", id)
				return mockArticle, nil
			},
		}
		fileRepo := &mockFileRepo{
			GetFileFunc: func(ctx context.Context, path string) ([]byte, error) {
				return mockContent, nil
			},
		}

		uc := NewArticleUsecase(articleRepo, fileRepo)

		// Act
		result, err := uc.GetArticle(ctx, "1")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "1", result.ID)
		assert.Equal(t, "Test Title", result.Title)
		assert.Equal(t, "Test Content", result.Content)
	})

	t.Run("ArticleRepo Error", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("database error")
		articleRepo := &mockArticleRepo{
			GetArticleFunc: func(ctx context.Context, id string) (*data.Article, error) {
				return nil, expectedErr
			},
		}
		fileRepo := &mockFileRepo{} // Not called, can be empty

		uc := NewArticleUsecase(articleRepo, fileRepo)

		// Act
		result, err := uc.GetArticle(ctx, "1")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestArticleUsecase_ListArticles(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockArticles := []*data.Article{
			{ID: "1", Title: "Article 1"},
			{ID: "2", Title: "Article 2"},
		}
		articleRepo := &mockArticleRepo{
			ListArticlesFunc: func(ctx context.Context, tag string, limit, offset int) ([]*data.Article, error) {
				return mockArticles, nil
			},
		}
		fileRepo := &mockFileRepo{} // Not called

		uc := NewArticleUsecase(articleRepo, fileRepo)

		// Act
		results, err := uc.ListArticles(ctx, "", 10, 0)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, 2)
		assert.Equal(t, "Article 1", results[0].Title)
	})
}
