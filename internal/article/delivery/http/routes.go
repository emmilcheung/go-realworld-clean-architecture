package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/article"
)

func ArticlesRouteRegister(router *gin.RouterGroup, h article.Handlers) {
	router.POST("/", h.ArticleCreate())
	router.PUT("/:slug", h.ArticleUpdate())
	router.DELETE("/:slug", h.ArticleDelete())
	router.POST("/:slug/favorite", h.ArticleFavorite())
	router.DELETE("/:slug/favorite", h.ArticleUnfavorite())
	router.POST("/:slug/comments", h.ArticleCommentCreate())
	router.DELETE("/:slug/comments/:id", h.ArticleCommentDelete())
}

func ArticlesAnonymousRouteRegister(router *gin.RouterGroup, h article.Handlers) {
	router.GET("/", h.ArticleList())
	router.GET("/feed", h.ArticleFeed())
	router.GET("/:slug", h.ArticleRetrieve())
	router.GET("/:slug/comments", h.ArticleCommentList())
}

func TagsAnonymousRouteRegister(router *gin.RouterGroup, h article.Handlers) {
	router.GET("/", h.TagList())
}
