package parserproxy
import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	parserpb "github.com/cosmtrek/blog-platform/api/parser/v1"
	proxypb "github.com/cosmtrek/blog-platform/api/proxy/v1"
)
type MockMarkdownClient struct { mock.Mock }
func (m *MockMarkdownClient) Parse(ctx context.Context, in *parserpb.ParseRequest, opts ...grpc.CallOption) (*parserpb.ParseResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*parserpb.ParseResponse), args.Error(1)
}
type MockTypstClient struct { mock.Mock }
func (m *MockTypstClient) Parse(ctx context.Context, in *parserpb.ParseRequest, opts ...grpc.CallOption) (*parserpb.ParseResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*parserpb.ParseResponse), args.Error(1)
}
func TestParserProxy_Parse(t *testing.T) {
	mockMarkdown := new(MockMarkdownClient)
	mockTypst := new(MockTypstClient)
	service := &Service{ markdownClient: mockMarkdown, typstClient: mockTypst }
	ctx := context.Background()
	testContent := []byte("test content")

	t.Run("markdown", func(t *testing.T) {
		req := &proxypb.ProxyParseRequest{Format: "markdown", Content: testContent}
		expectedResp := &parserpb.ParseResponse{Content: []byte("markdown html")}
		mockMarkdown.On("Parse", ctx, &parserpb.ParseRequest{Content: testContent}).Return(expectedResp, nil).Once()
		resp, err := service.Parse(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, expectedResp.Content, resp.Content)
		mockMarkdown.AssertExpectations(t)
	})
	t.Run("typst", func(t *testing.T) {
		req := &proxypb.ProxyParseRequest{Format: "typst", Content: testContent}
		expectedResp := &parserpb.ParseResponse{Content: []byte("typst html")}
		mockTypst.On("Parse", ctx, &parserpb.ParseRequest{Content: testContent}).Return(expectedResp, nil).Once()
		resp, err := service.Parse(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, expectedResp.Content, resp.Content)
		mockTypst.AssertExpectations(t)
	})
}
