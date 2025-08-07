package service

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

type Resolver struct{}

// NewExecutableSchema creates a new GraphQL schema.
func NewExecutableSchema(cfg Config) *handler.Server {
	return handler.NewDefaultServer(nil) // This is a placeholder
}

type Config struct {
	Resolvers *Resolver
}

func (c *Config) Complexity(typeName, field string, childComplexity int, args map[string]interface{}) (int, bool) {
	return 0, false
}

func (c *Config) Directives() DireciveRoot {
	return DireciveRoot{}
}

type DireciveRoot struct{}

func (c *Config) Query() QueryResolver {
	return &queryResolver{c.Resolvers}
}

type QueryResolver interface {
	Article(ctx context.Context, id string) (*Article, error)
	Articles(ctx context.Context, tag *string, limit *int) ([]*Article, error)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Article(ctx context.Context, id string) (*Article, error) {
	return r.Resolver.Article(ctx, id)
}
func (r *queryResolver) Articles(ctx context.Context, tag *string, limit *int) ([]*Article, error) {
	return r.Resolver.Articles(ctx, tag, limit)
}

// ---- Playground Handler ----

func PlaygroundHandler(w http.ResponseWriter, r *http.Request) {
	playground.Handler("GraphQL playground", "/graphql").ServeHTTP(w, r)
}
