package models

import (
	"time"
)

type Article struct {
	ID          uint   `gorm:"primaryKey"`
	Slug        string `gorm:"unique_index"`
	Title       string
	Description string `gorm:"size:2048"`
	Body        string `gorm:"size:2048"`
	Author      ArticleUser
	AuthorID    uint
	Tags        []Tag     `gorm:"many2many:article_tags;"`
	Comments    []Comment `gorm:"ForeignKey:ArticleID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index" json:"deleted_at"`
}

func (e *Article) TableName() string {
	return "article_models"
}

type ArticleUser struct {
	ID        uint `gorm:"primaryKey"`
	User      User
	UserID    uint
	Articles  []Article  `gorm:"ForeignKey:AuthorID"`
	Favorites []Favorite `gorm:"ForeignKey:FavoriteByID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (e *ArticleUser) TableName() string {
	return "article_user_models"
}

type Favorite struct {
	ID           uint `gorm:"primaryKey"`
	Favorite     Article
	FavoriteID   uint
	FavoriteBy   ArticleUser
	FavoriteByID uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index" json:"deleted_at"`
}

func (e *Favorite) TableName() string {
	return "favorite_models"
}

type Tag struct {
	ID        uint      `gorm:"primaryKey"`
	Tag       string    `gorm:"unique_index"`
	Articles  []Article `gorm:"many2many:article_tags;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (e *Tag) TableName() string {
	return "tag_models"
}

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	Article   Article
	ArticleID uint
	Author    ArticleUser
	AuthorID  uint
	Body      string `gorm:"size:2048"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (e *Comment) TableName() string {
	return "comment_models"
}
