package biz

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFileRepo is a mock of FileRepo interface
type MockFileRepo struct {
	mock.Mock
}

func (m *MockFileRepo) Upload(ctx context.Context, file *File) (string, error) {
	args := m.Called(ctx, file)
	return args.String(0), args.Error(1)
}

func (m *MockFileRepo) Download(ctx context.Context, path string) (*File, error) {
	args := m.Called(ctx, path)
	// Need to handle the case where the first return value is nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*File), args.Error(1)
}

func (m *MockFileRepo) PublishUploadEvent(ctx context.Context, filePath string, format string) error {
	args := m.Called(ctx, filePath, format)
	return args.Error(0)
}

func TestFileUsecase_UploadFile(t *testing.T) {
	mockRepo := new(MockFileRepo)
	logger := log.NewStdLogger(io.Discard)
	usecase := NewFileUsecase(mockRepo, logger)

	ctx := context.Background()
	file := &File{
		Name:    "test.md",
		Content: strings.NewReader("hello world"),
		Size:    11,
	}
	format := "markdown"
	filePath := "articles/test.md"

	// Setup expectations
	mockRepo.On("Upload", ctx, file).Return(filePath, nil).Once()
	mockRepo.On("PublishUploadEvent", ctx, filePath, format).Return(nil).Once()

	// Call the method
	resultPath, err := usecase.UploadFile(ctx, file, format)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, filePath, resultPath)
	mockRepo.AssertExpectations(t)
}

func TestFileUsecase_UploadFile_UploadFails(t *testing.T) {
	mockRepo := new(MockFileRepo)
	logger := log.NewStdLogger(io.Discard)
	usecase := NewFileUsecase(mockRepo, logger)

	ctx := context.Background()
	file := &File{
		Name:    "test.md",
		Content: strings.NewReader("hello world"),
		Size:    11,
	}
	format := "markdown"
	uploadError := errors.New("upload failed")

	// Setup expectations
	mockRepo.On("Upload", ctx, file).Return("", uploadError).Once()

	// Call the method
	_, err := usecase.UploadFile(ctx, file, format)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, uploadError, err)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "PublishUploadEvent", mock.Anything, mock.Anything, mock.Anything)
}
