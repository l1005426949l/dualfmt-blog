package data

import (
	"context"
	"database/sql"
	"my-project/metadata/biz"
	"my-project/metadata/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/marcboeker/go-duckdb/v2"
)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(logger)
	sqlDB, err := sql.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		logHelper.Errorf("failed opening connection to %s: %v", c.Database.Driver, err)
		return nil, nil, err
	}

	db, err := gorm.Open(sqlite.Dialector{
		Conn: sqlDB,
	}, &gorm.Config{})

	if err != nil {
		logHelper.Errorf("failed opening gorm connection: %v", err)
		return nil, nil, err
	}

	if err := db.AutoMigrate(&Article{}); err != nil {
		logHelper.Errorf("failed auto migrating tables: %v", err)
		return nil, nil, err
	}

	cleanup := func() {
		logHelper.Info("closing the data resources")
		sqlDB.Close()
	}
	return &Data{db: db}, cleanup, nil
}

type articleRepo struct {
	data *Data
	log  *log.Helper
}

func NewArticleRepo(data *Data, logger log.Logger) biz.ArticleRepo {
	return &articleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *articleRepo) Get(ctx context.Context, id string) (*biz.Article, error) {
	var article Article
	if err := r.data.db.WithContext(ctx).Where("article_id = ?", id).First(&article).Error; err != nil {
		return nil, err
	}
	return &biz.Article{
		ID:      article.ArticleID,
		Title:   article.Title,
		Content: article.Content,
		Version: article.Version,
		Format:  article.Format,
	}, nil
}

func (r *articleRepo) List(ctx context.Context, tag string, limit int) ([]*biz.Article, error) {
	var articles []Article
	// DuckDB does not support tags, so we will ignore the tag for now.
	if limit > 0 {
		if err := r.data.db.WithContext(ctx).Limit(limit).Find(&articles).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.data.db.WithContext(ctx).Find(&articles).Error; err != nil {
			return nil, err
		}
	}

	var bizArticles []*biz.Article
	for _, article := range articles {
		bizArticles = append(bizArticles, &biz.Article{
			ID:      article.ArticleID,
			Title:   article.Title,
			Content: article.Content,
			Version: article.Version,
			Format:  article.Format,
		})
	}
	return bizArticles, nil
}

var ProviderSet = wire.NewSet(NewData, NewArticleRepo)
