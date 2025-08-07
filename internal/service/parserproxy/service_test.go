package parserproxy

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	parserpb "github.com/cosmtrek/blog-platform/api/parser/v1"
	proxypb "github.com/cosmtrek/blog-platform/api/proxy/v1"
)

// MockMarkdownClient is a mock of MarkdownParserClient
type MockMarkdownClient struct {
	mock.Mock
}

func (m *MockMarkdownClient) Parse(ctx context.Context, in *parserpb.ParseRequest, opts ...grpc.CallOption) (*parserpb.ParseResponse, error) {
	args := m.Called(ctx, in)
	// Handle the case where the response is nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*parserpb.ParseResponse), args.Error(1)
}

// MockTypstClient is a mock of TypstParserClient
type MockTypstClient struct {
	mock.Mock
}

func (m *MockTypstClient) Parse(ctx context.Context, in *parserpb.ParseRequest, opts ...grpc.CallOption) (*parserpb.ParseResponse, error) {
	args := m.Called(ctx, in)
	// Handle the case where the response is nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*parserpb.ParseResponse), args.Error(1)
}


func TestParserProxy_Parse(t *testing.T) {
	// Setup
	mockMarkdown := new(MockMarkdownClient)
	mockTypst := new(MockTypstClient)

	service := &Service{
		markdownClient: mockMarkdown,
		typstClient:    mockTypst,
	}

	ctx := context.Background()
	testContent := []byte("test content")

	// Test case 1: Markdown
	t.Run("should route to markdown parser", func(t *testing.T) {
		req := &proxypb.ProxyParseRequest{Format: "markdown", Content: testContent}
		expectedResp := &parserpb.ParseResponse{Content: []byte("markdown html")}

		mockMarkdown.On("Parse", ctx, &parserpb.ParseRequest{Content: testContent}).Return(expectedResp, nil).Once()
		mockTypst.AssertNotCalled(t, "Parse")

		resp, err := service.Parse(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResp.Content, resp.Content)
		mockMarkdown.AssertExpectations(t)
	})

	// Test case 2: Typst
	t.Run("should route to typst parser", func(t *testing.T) {
		req := &proxypb.ProxyParseRequest{Format: "typst", Content: testContent}
		expectedResp := &parserpb.ParseResponse{Content: []byte("typst html")}

		mockTypst.On("Parse", ctx, &parserpb.ParseRequest{Content: testContent}).Return(expectedResp, nil).Once()
		mockMarkdown.AssertNotCalled(t, "Parse")

		resp, err := service.Parse(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResp.Content, resp.Content)
		mockTypst.AssertExpectations(t)
	})

	// Test case 3: Unsupported format
	t.Run("should return error for unsupported format", func(t *testing.T) {
		req := &proxypb.ProxyParseRequest{Format: "invalid", Content: testContent}

		resp, err := service.Parse(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockMarkdown.AssertNotCalled(t, "Parse")
		mockTypst.AssertNotCalled(t, "Parse")
	})

    // Test case 4: Markdown parser error
	t.Run("should return error when markdown parser fails", func(t *testing.T) {
		req := &proxypb.ProxyParseRequest{Format: "markdown", Content: testContent}
        expectedErr := errors.New("markdown error")

		mockMarkdown.On("Parse", ctx, &parserpb.ParseRequest{Content: testContent}).Return(nil, expectedErr).Once()
		mockTypst.AssertNotCalled(t, "Parse")

		resp, err := service.Parse(ctx, req)

		assert.Error(t, err)
        assert.Contains(t, err.Error(), "markdown parser failed")
		assert.Nil(t, resp)
		mockMarkdown.AssertExpectations(t)
	})
}
