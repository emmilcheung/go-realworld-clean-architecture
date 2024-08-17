package article

import (
	"context"

	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
)

type Repository interface {
	GetArticleUser(c context.Context, userID uint) models.ArticleUser
	FindManyArticle(c context.Context, tag, author, favorited string, limit, offset int) ([]models.Article, int, error)
	FindOneArticle(c context.Context, condition interface{}) (models.Article, error)
	GetArticleFeed(c context.Context, userId uint, limit, offset int) ([]models.Article, int, error)
	SaveOne(ctx context.Context, data interface{}) error
	Update(c context.Context, data *models.Article) error
	DeleteArticleModel(c context.Context, condition interface{}) error
	UpsertTags(ctx context.Context, tags []string) ([]models.Tag, error)
	ArticleFavoritesCount(c context.Context, articleId uint) uint
	IsArticleFavoriteBy(c context.Context, userId uint, articleId uint) bool
	SetFavorite(ctx context.Context, articleId, userId uint) error
	RemoveFavorite(ctx context.Context, articleId, userId uint) error
	DeleteComment(ctx context.Context, condition interface{}) error
	GetArticleComments(ctx context.Context, article models.Article) ([]models.Comment, error)
	GetTags() ([]models.Tag, error)
}
