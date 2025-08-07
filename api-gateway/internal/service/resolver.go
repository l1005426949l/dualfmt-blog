package service

import (
	"context"
	"api-gateway/internal/biz"
	"api-gateway/internal/data"
)

// The resolver now holds a reference to our business logic use cases.
// It is defined in generated.go, so we modify its dependencies here.
type Resolver struct{
	uc *biz.ArticleUsecase
}

// NewResolver creates a new resolver with its dependencies.
func NewResolver(uc *biz.ArticleUsecase) *Resolver {
	return &Resolver{uc: uc}
}

// Article resolves the Article query by calling the business logic layer.
func (r *Resolver) Article(ctx context.Context, id string) (*Article, error) {
	bizArticle, err := r.uc.GetArticle(ctx, id)
	if err != nil {
		return nil, err
	}
	// Map from biz/data model to GraphQL model
	return mapBizArticleToGQL(bizArticle), nil
}

// Articles resolves the Articles query by calling the business logic layer.
func (r *Resolver) Articles(ctx context.Context, tag *string, limit *int) ([]*Article, error) {
	// Set defaults for limit
	l := 10
	if limit != nil {
		l = *limit
	}

	var t string
	if tag != nil {
		t = *tag
	}

	bizArticles, err := r.uc.ListArticles(ctx, t, l, 0)
	if err != nil {
		return nil, err
	}

	// Map from biz/data model to GraphQL model
	gqlArticles := make([]*Article, len(bizArticles))
	for i, a := range bizArticles {
		gqlArticles[i] = mapBizArticleToGQL(a)
	}
	return gqlArticles, nil
}

// mapBizArticleToGQL maps the internal article model to the GraphQL article model.
func mapBizArticleToGQL(a *data.Article) *Article {
	if a == nil {
		return nil
	}
	return &Article{
		ID:      a.ID,
		Title:   a.Title,
		Content: a.Content,
		Version: a.Version,
		Format:  a.Format,
	}
}
