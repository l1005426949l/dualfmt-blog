package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

type Article struct {
	ID      string
	Title   string
	Content string
	Version string
	Format  string
}

type ArticleUsecase struct {
	repo ArticleRepo
	log  *log.Helper
}

func NewArticleUsecase(repo ArticleRepo, logger log.Logger) *ArticleUsecase {
	return &ArticleUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *ArticleUsecase) Get(ctx context.Context, id string) (*Article, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *ArticleUsecase) List(ctx context.Context, tag string, limit int) ([]*Article, error) {
	return uc.repo.List(ctx, tag, limit)
}

var ProviderSet = wire.NewSet(NewArticleUsecase)
