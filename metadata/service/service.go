package service

import (
	"context"

	pb "my-project/api/metadata/v1"
	"my-project/metadata/biz"
	"github.com/google/wire"
)

type MetadataService struct {
	pb.UnimplementedMetadataServiceServer
	uc *biz.ArticleUsecase
}

func NewMetadataService(uc *biz.ArticleUsecase) *MetadataService {
	return &MetadataService{uc: uc}
}

func (s *MetadataService) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleReply, error) {
	article, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetArticleReply{
		Article: &pb.Article{
			Id:      article.ID,
			Title:   article.Title,
			Content: article.Content,
			Version: article.Version,
			Format:  article.Format,
		},
	}, nil
}

func (s *MetadataService) ListArticles(ctx context.Context, req *pb.ListArticlesRequest) (*pb.ListArticlesReply, error) {
	articles, err := s.uc.List(ctx, req.Tag, int(req.Limit))
	if err != nil {
		return nil, err
	}
	var pbArticles []*pb.Article
	for _, article := range articles {
		pbArticles = append(pbArticles, &pb.Article{
			Id:      article.ID,
			Title:   article.Title,
			Content: article.Content,
			Version: article.Version,
			Format:  article.Format,
		})
	}
	return &pb.ListArticlesReply{Articles: pbArticles}, nil
}

var ProviderSet = wire.NewSet(NewMetadataService)
