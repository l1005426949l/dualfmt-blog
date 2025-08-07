package parserproxy
import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	parserpb "github.com/cosmtrek/blog-platform/api/parser/v1"
	proxypb "github.com/cosmtrek/blog-platform/api/proxy/v1"
)
type Service struct {
	proxypb.UnimplementedParserProxyServer
	markdownClient parserpb.MarkdownParserClient
	typstClient    parserpb.TypstParserClient
}
func New(markdownAddr, typstAddr string) (*Service, error) {
	markdownConn, err := grpc.Dial(markdownAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil { return nil, err }
	typstConn, err := grpc.Dial(typstAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil { return nil, err }
	return &Service{
		markdownClient: parserpb.NewMarkdownParserClient(markdownConn),
		typstClient:    parserpb.NewTypstParserClient(typstConn),
	}, nil
}
func (s *Service) Parse(ctx context.Context, req *proxypb.ProxyParseRequest) (*proxypb.ProxyParseResponse, error) {
	switch req.Format {
	case "markdown":
		resp, err := s.markdownClient.Parse(ctx, &parserpb.ParseRequest{Content: req.Content})
		if err != nil { return nil, fmt.Errorf("markdown parser failed: %w", err) }
		return &proxypb.ProxyParseResponse{Content: resp.Content}, nil
	case "typst":
		resp, err := s.typstClient.Parse(ctx, &parserpb.ParseRequest{Content: req.Content})
		if err != nil { return nil, fmt.Errorf("typst parser failed: %w", err) }
		return &proxypb.ProxyParseResponse{Content: resp.Content}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", req.Format)
	}
}
