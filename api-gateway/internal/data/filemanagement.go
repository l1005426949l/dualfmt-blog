package data

import (
	"context"
	"fmt"
)

// mockFileRepo is a mock implementation of FileRepo.
type mockFileRepo struct{}

// UploadFile simulates uploading a file and returns a fake path.
func (r *mockFileRepo) UploadFile(ctx context.Context, filename string, content []byte) (string, error) {
	fmt.Printf("Simulating upload for file: %s, size: %d bytes\n", filename, len(content))
	// Return a predictable, fake storage path
	return "mock/storage/path/for/" + filename, nil
}

// GetFile simulates fetching a file's content.
func (r *mockFileRepo) GetFile(ctx context.Context, path string) ([]byte, error) {
	fmt.Printf("Simulating fetch for file at path: %s\n", path)
	// Return some mock content based on the path
	if path == "mock/storage/path/for/article-1.md" {
		return []byte("# Hello World!\n\nThis is the content from the mock file repo."), nil
	}
	return []byte("Mock file content for path: " + path), nil
}
