package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/article"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/httpErrors"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/locker"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/utils"
	"github.com/opentracing/opentracing-go"
)

type articleHandlers struct {
	articleRepo article.Repository
	articleUc   article.UseCase
	locker      locker.Locker
}

func NewArticleHandlers(articleRepo article.Repository, articleUc article.UseCase, locker locker.Locker) article.Handlers {
	return &articleHandlers{articleRepo, articleUc, locker}
}

func (h articleHandlers) ArticleList() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "article.ArticleList")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()
		tag := c.Query("tag")
		author := c.Query("author")
		favorited := c.Query("favorited")
		pagination, err := utils.GetPaginationFromCtx(c)
		articleModels, modelCount, err := h.articleUc.GetArticles(ctx, tag, author, favorited, pagination)
		if err != nil {
			c.JSON(http.StatusNotFound, httpErrors.NewError("articles", errors.New("Invalid param")))
			return
		}

		serializer := ArticlesSerializer{ctx, h.articleRepo, articleModels}
		c.JSON(http.StatusOK, gin.H{"articles": serializer.Response(), "articlesCount": modelCount})

	}
}

func (h articleHandlers) ArticleRetrieve() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "article.ArticleRetrieve")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		slug := c.Param("slug")
		articleModel, err := h.articleUc.GetArticle(ctx, slug)
		if err != nil {
			c.JSON(http.StatusNotFound, httpErrors.NewError("articles", errors.New("Invalid slug")))
			return
		}
		serializer := ArticleSerializer{ctx, h.articleRepo, articleModel}
		c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
	}
}

func (h articleHandlers) ArticleFeed() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "article.ArticleFeed")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()
		pagination, err := utils.GetPaginationFromCtx(c)
		myUserModel := c.MustGet("my_user_model").(models.User)
		if myUserModel.ID == 0 {
			c.AbortWithError(http.StatusUnauthorized, errors.New("{error : \"Require auth!\"}"))
			return
		}
		articleModels, modelCount, err := h.articleUc.GetFeeds(ctx, myUserModel, pagination)
		if err != nil {
			c.JSON(http.StatusNotFound, httpErrors.NewError("articles", errors.New("Invalid param")))
			return
		}
		serializer := ArticlesSerializer{ctx, h.articleRepo, articleModels}
		c.JSON(http.StatusOK, gin.H{"articles": serializer.Response(), "articlesCount": modelCount})
	}
}

func (h articleHandlers) ArticleCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "article.ArticleCreate")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		myUserModel := c.MustGet("my_user_model").(models.User)
		articleModelValidator := NewArticleModelValidator()
		if err := articleModelValidator.Verify(c); err != nil {
			c.JSON(http.StatusUnprocessableEntity, httpErrors.NewValidatorError(err))
			return
		}

		articleSlug := slug.Make(articleModelValidator.Article.Title)
		lockKey := fmt.Sprintf("article:slug-%s", articleSlug)
		lock := h.locker.ObtainLock(ctx, lockKey)
		defer lock.Release(ctx)

		articleModel, err := h.articleUc.CreateArticle(ctx, &models.Article{
			Slug:        articleSlug,
			Title:       articleModelValidator.Article.Title,
			Description: articleModelValidator.Article.Description,
			Body:        articleModelValidator.Article.Body,
			AuthorID:    myUserModel.ID,
		}, articleModelValidator.Article.Tags)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, httpErrors.NewValidatorError(err))
			return
		}

		serializer := ArticleSerializer{ctx, h.articleRepo, *articleModel}
		c.JSON(http.StatusCreated, gin.H{"article": serializer.Response()})
	}
}

func (h articleHandlers) ArticleUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "article.ArticleUpdate")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()
		articleSlug := c.Param("slug")
		articleModelValidator := NewArticlePartialModelValidator()
		if err := articleModelValidator.Verify(c); err != nil {
			c.JSON(http.StatusUnprocessableEntity, httpErrors.NewValidatorError(err))
			return
		}

		myUserModel := c.MustGet("my_user_model").(models.User)

		// apply lock to new slug to prevent collision
		lockKey := fmt.Sprintf("article:slug-%s", articleSlug)
		lock := h.locker.ObtainLock(ctx, lockKey)
		defer lock.Release(ctx)

		articleModel, err := h.articleUc.UpdateArticle(ctx, articleSlug, &models.Article{
			Slug:        slug.Make(articleModelValidator.Article.Title),
			Title:       articleModelValidator.Article.Title,
			Description: articleModelValidator.Article.Description,
			Body:        articleModelValidator.Article.Body,
			AuthorID:    myUserModel.ID,
		}, articleModelValidator.Article.Tags)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, httpErrors.NewError("database", err))
			return
		}
		serializer := ArticleSerializer{ctx, h.articleRepo, *articleModel}
		c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
	}
}

func (h articleHandlers) ArticleDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "article.ArticleDelete")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		slug := c.Param("slug")
		lockKey := fmt.Sprintf("article:slug-%s", slug)
		lock := h.locker.ObtainLock(ctx, lockKey)
		defer lock.Release(ctx)

		if err := h.articleUc.DeleteArticle(ctx, slug); err != nil {
			c.JSON(http.StatusNotFound, httpErrors.NewError("articles", errors.New("Invalid slug")))
			return
		}
		c.JSON(http.StatusOK, gin.H{"article": "Delete success"})
	}
}

func (h articleHandlers) ArticleFavorite() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "article.ArticleFavorite")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		articleSlug := c.Param("slug")
		myUserModel := c.MustGet("my_user_model").(models.User)
		articleModel, err := h.articleUc.CreateFavorite(ctx, articleSlug, myUserModel.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, httpErrors.NewError("articles", errors.New("Invalid slug")))
			return
		}
		serializer := ArticleSerializer{ctx, h.articleRepo, *articleModel}
		c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
	}
}

func (h articleHandlers) ArticleUnfavorite() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "article.ArticleUnfavorite")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		slug := c.Param("slug")
		myUserModel := c.MustGet("my_user_model").(models.User)
		articleModel, err := h.articleUc.DeleteFavorite(ctx, slug, myUserModel.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, httpErrors.NewError("articles", err))
			return
		}
		serializer := ArticleSerializer{ctx, h.articleRepo, *articleModel}
		c.JSON(http.StatusOK, gin.H{"article": serializer.Response()})
	}
}

func (h articleHandlers) ArticleCommentList() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "article.ArticleCommentList")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		slug := c.Param("slug")
		comments, err := h.articleUc.GetCommentsByArticle(ctx, slug)
		if err != nil {
			c.JSON(http.StatusNotFound, httpErrors.NewError("comments", err))
		}
		serializer := CommentsSerializer{ctx, comments}
		c.JSON(http.StatusOK, gin.H{"comments": serializer.Response()})
	}
}

func (h articleHandlers) ArticleCommentCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "article.ArticleCommentCreate")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		slug := c.Param("slug")
		myUserModel := c.MustGet("my_user_model").(models.User)

		commentModelValidator := NewCommentModelValidator()
		if err := commentModelValidator.Verify(c); err != nil {
			c.JSON(http.StatusUnprocessableEntity, httpErrors.NewValidatorError(err))
			return
		}

		commentModel, err := h.articleUc.CreateComment(ctx, slug, myUserModel.ID, &models.Comment{
			Body: commentModelValidator.Comment.Body,
		})

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		serializer := CommentSerializer{ctx, *commentModel}
		c.JSON(http.StatusCreated, gin.H{"comment": serializer.Response()})
	}
}

func (h articleHandlers) ArticleCommentDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "article.ArticleCommentDelete")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
		myUserModel := c.MustGet("my_user_model").(models.User)
		id := uint(id64)
		if err != nil {
			c.JSON(http.StatusNotFound, httpErrors.NewError("comment", errors.New("Invalid id")))
			return
		}
		err = h.articleUc.DeleteComment(ctx, myUserModel.ID, []uint{id})
		if err != nil {
			c.JSON(http.StatusNotFound, httpErrors.NewError("comment", err))
			return
		}
		c.JSON(http.StatusOK, gin.H{"comment": "Delete success"})
	}
}

func (h articleHandlers) TagList() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "article.TagList")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		tagModels, err := h.articleRepo.GetTags()
		if err != nil {
			c.JSON(http.StatusNotFound, httpErrors.NewError("articles", err))
			return
		}
		serializer := TagsSerializer{ctx, tagModels}
		c.JSON(http.StatusOK, gin.H{"tags": serializer.Response()})
	}
}
