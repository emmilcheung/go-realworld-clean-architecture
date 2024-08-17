package article

import (
	"context"

	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/utils"
)

type UseCase interface {
	GetArticleUser(ctx context.Context, userID uint) models.ArticleUser
	GetArticles(ctx context.Context, tag, author, favorited string, pagination *utils.PaginationQuery) ([]models.Article, int, error)
	GetFeeds(ctx context.Context, user models.User, pagination *utils.PaginationQuery) ([]models.Article, int, error)
	GetArticle(ctx context.Context, slug string) (models.Article, error)
	CreateArticle(ctx context.Context, articleModel *models.Article, tags []string) (*models.Article, error)
	UpdateArticle(ctx context.Context, slug string, articleModel *models.Article, tags []string) (*models.Article, error)
	DeleteArticle(ctx context.Context, slug string) error
	CreateFavorite(ctx context.Context, slug string, userID uint) (*models.Article, error)
	DeleteFavorite(ctx context.Context, slug string, userID uint) (*models.Article, error)
	GetCommentsByArticle(ctx context.Context, slug string) ([]models.Comment, error)
	CreateComment(ctx context.Context, slug string, userID uint, comment *models.Comment) (*models.Comment, error)
	DeleteComment(ctx context.Context, userID uint, commentIDs []uint) error
}
