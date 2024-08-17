package usecase

import (
	"context"
	"errors"

	"github.com/gothinkster/golang-gin-realworld-example-app/internal/article"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/utils"
	"github.com/opentracing/opentracing-go"
)

// Comments UseCase
type articleUC struct {
	// cfg      *config.Config
	// logger   logger.Logger
	articleRepo article.Repository
}

// Comments UseCase constructor
func NewArticleUseCase(articleRepo article.Repository) article.UseCase {
	return &articleUC{articleRepo: articleRepo}
}

func (uc *articleUC) GetArticleUser(ctx context.Context, userID uint) models.ArticleUser {
	return uc.articleRepo.GetArticleUser(ctx, userID)
}

func (uc *articleUC) GetArticles(ctx context.Context, tag, author, favorited string, pagination *utils.PaginationQuery) ([]models.Article, int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "article.usecase.GetArticles")
	defer span.Finish()
	return uc.articleRepo.FindManyArticle(ctx, tag, author, favorited, pagination.Limit, pagination.Offset)
}

func (uc *articleUC) GetFeeds(ctx context.Context, user models.User, pagination *utils.PaginationQuery) ([]models.Article, int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "article.usecase.GetFeeds")
	defer span.Finish()

	articleUserModel := uc.articleRepo.GetArticleUser(ctx, user.ID)
	return uc.articleRepo.GetArticleFeed(ctx, articleUserModel.UserID, pagination.Limit, pagination.Offset)
}

func (uc *articleUC) GetArticle(ctx context.Context, slug string) (models.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "article.usecase.GetArticle")
	defer span.Finish()
	return uc.articleRepo.FindOneArticle(ctx, &models.Article{Slug: slug})
}

func (uc *articleUC) CreateArticle(ctx context.Context, articleModel *models.Article, tags []string) (*models.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "article.usecase.CreateArticle")
	defer span.Finish()
	articleModel.Author = uc.articleRepo.GetArticleUser(ctx, articleModel.AuthorID)
	tagModels, _ := uc.articleRepo.UpsertTags(ctx, tags)
	articleModel.Tags = tagModels

	err := uc.articleRepo.SaveOne(ctx, &articleModel)
	return articleModel, err
}

func (uc *articleUC) UpdateArticle(ctx context.Context, slug string, updateArticle *models.Article, tags []string) (*models.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "article.usecase.UpdateArticle")
	defer span.Finish()

	articleModel, err := uc.articleRepo.FindOneArticle(ctx, &models.Article{Slug: slug})
	if err != nil {
		return nil, err
	}
	if updateArticle.Title != "" {
		articleModel.Title = updateArticle.Title
	}
	if updateArticle.Description != "" {
		articleModel.Description = updateArticle.Description
	}
	if updateArticle.Body != "" {
		articleModel.Body = updateArticle.Body
	}

	articleModel.Author = uc.articleRepo.GetArticleUser(ctx, articleModel.AuthorID)
	tagModels, _ := uc.articleRepo.UpsertTags(ctx, tags)
	articleModel.Tags = tagModels

	err = uc.articleRepo.Update(ctx, &articleModel)
	return &articleModel, err
}

func (uc *articleUC) DeleteArticle(ctx context.Context, slug string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "article.usecase.DeleteArticle")
	defer span.Finish()

	_, err := uc.articleRepo.FindOneArticle(ctx, &models.Article{Slug: slug})
	if err != nil {
		return err
	}

	return uc.articleRepo.DeleteArticleModel(ctx, &models.Article{Slug: slug})
}

func (uc *articleUC) CreateFavorite(ctx context.Context, slug string, userID uint) (*models.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "article.usecase.CreateFavorite")
	defer span.Finish()

	articleModel, err := uc.articleRepo.FindOneArticle(ctx, &models.Article{Slug: slug})
	if err != nil {
		return nil, err
	}
	uc.articleRepo.SetFavorite(ctx, articleModel.ID, uc.articleRepo.GetArticleUser(ctx, userID).ID)

	return &articleModel, err
}
func (uc *articleUC) DeleteFavorite(ctx context.Context, slug string, userID uint) (*models.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "article.usecase.CreateFavorite")
	defer span.Finish()

	articleModel, err := uc.articleRepo.FindOneArticle(ctx, &models.Article{Slug: slug})
	if err != nil {
		return nil, err
	}
	uc.articleRepo.RemoveFavorite(ctx, articleModel.ID, uc.articleRepo.GetArticleUser(ctx, userID).ID)

	return &articleModel, err
}

func (uc *articleUC) GetCommentsByArticle(ctx context.Context, slug string) ([]models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "article.usecase.GetCommentsByArticle")
	defer span.Finish()

	articleModel, err := uc.articleRepo.FindOneArticle(ctx, (&models.Article{Slug: slug}))
	if err != nil {
		return nil, errors.New("Invalid slug")
	}
	comments, err := uc.articleRepo.GetArticleComments(ctx, articleModel)
	if err != nil {
		return nil, errors.New("Database error")
	}
	return comments, nil
}

func (uc *articleUC) CreateComment(ctx context.Context, slug string, userID uint, comment *models.Comment) (*models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "article.usecase.CreateComment")
	defer span.Finish()

	articleModel, err := uc.articleRepo.FindOneArticle(ctx, &models.Article{Slug: slug})
	if err != nil {
		return nil, errors.New("Invalid slug")
	}

	comment.Article = articleModel
	comment.Author = uc.articleRepo.GetArticleUser(ctx, userID)

	err = uc.articleRepo.SaveOne(ctx, comment)
	return comment, err
}

func (uc *articleUC) DeleteComment(ctx context.Context, userID uint, commentIDs []uint) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "article.usecase.DeleteComment")
	defer span.Finish()

	return uc.articleRepo.DeleteComment(ctx, commentIDs)
}

func (uc *articleUC) GetTags(ctx context.Context) ([]models.Tag, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "article.usecase.GetTags")
	defer span.Finish()
	return uc.articleRepo.GetTags()
}
