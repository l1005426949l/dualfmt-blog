package biz

import (
	"context"
)

type ArticleRepo interface {
	Get(ctx context.Context, id string) (*Article, error)
	List(ctx context.Context, tag string, limit int) ([]*Article, error)
}
