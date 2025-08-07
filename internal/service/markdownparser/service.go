package markdownparser
import (
	"context"
	pb "github.com/cosmtrek/blog-platform/api/parser/v1"
)
type Service struct { pb.UnimplementedMarkdownParserServer }
func New() *Service { return &Service{} }
func (s *Service) Parse(ctx context.Context, req *pb.ParseRequest) (*pb.ParseResponse, error) {
	return &pb.ParseResponse{Content: []byte("<h1>Hello from MarkdownParser</h1>")}, nil
}
