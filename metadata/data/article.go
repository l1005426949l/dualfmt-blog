package data

import (
	"time"
)

type Article struct {
	ID        uint `gorm:"primarykey"`
	ArticleID string
	Title     string
	Content   string
	Version   string
	Format    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
