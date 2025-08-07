package main

import (
	"context"
	"log"
	"net/http"

	appMiddleware "api-gateway/internal/middleware"
	"api-gateway/internal/service"
	"api-gateway/internal/conf"
	"api-gateway/internal/data"
	"api-gateway/internal/biz"


	"os"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	// --- 1. Load Configuration ---
	bc, err := conf.LoadBootstrap()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// --- 2. Initialize Layers (Data -> Biz -> Service) ---
	// This is manual dependency injection. Kratos uses Wire for this.
	dataRepo, cleanup, err := data.NewData()
	if err != nil {
		log.Fatalf("failed to create data layer: %v", err)
	}
	defer cleanup()

	articleUsecase := biz.NewArticleUsecase(dataRepo.ArticleRepo, dataRepo.FileRepo)
	resolver := service.NewResolver(articleUsecase)

	// --- 3. Create GraphQL server ---
	// This is the manual setup that replaces the generated code.
	schemaBytes, err := os.ReadFile("./internal/service/schema.graphqls")
	if err != nil {
		log.Fatalf("failed to read schema: %v", err)
	}

	executableSchema := graphql.NewExecutableSchema(graphql.Config{
		Resolvers: resolver,
		Sources: []*ast.Source{
			{Name: "schema.graphqls", Input: string(schemaBytes)},
		},
	})
	srv := handler.NewDefaultServer(executableSchema)


	// --- 4. Create Kratos HTTP server with middleware ---
	var opts []kratosHttp.ServerOption
	if bc.Server.HTTP.Network != "" {
		opts = append(opts, kratosHttp.Network(bc.Server.HTTP.Network))
	}
	if bc.Server.HTTP.Addr != "" {
		opts = append(opts, kratosHttp.Address(bc.Server.HTTP.Addr))
	}
	if bc.Server.HTTP.Timeout != 0 {
		opts = append(opts, kratosHttp.Timeout(bc.Server.HTTP.Timeout))
	}
	// Add middleware
	opts = append(opts, kratosHttp.Middleware(
		recovery.Recovery(),
		ratelimit.Server(),
		circuitbreaker.Server(),
		appMiddleware.JWTAuthMiddleware(),
	))

	httpSrv := kratosHttp.NewServer(opts...)

	// Register the GraphQL endpoint
	httpSrv.Handle("/graphql", srv)

	// Register the GraphQL Playground for easy testing in a browser
	httpSrv.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))


	// --- 4. Create and run Kratos app ---
	app := kratos.New(
		kratos.Name("api-gateway"),
		kratos.Version("v1.0.0"),
		kratos.Server(
			httpSrv,
		),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
