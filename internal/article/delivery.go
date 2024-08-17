package article

import "github.com/gin-gonic/gin"

type Handlers interface {
	ArticleCreate() gin.HandlerFunc
	ArticleList() gin.HandlerFunc
	ArticleRetrieve() gin.HandlerFunc
	ArticleFeed() gin.HandlerFunc
	ArticleUpdate() gin.HandlerFunc
	ArticleDelete() gin.HandlerFunc
	ArticleFavorite() gin.HandlerFunc
	ArticleUnfavorite() gin.HandlerFunc
	ArticleCommentCreate() gin.HandlerFunc
	ArticleCommentDelete() gin.HandlerFunc
	ArticleCommentList() gin.HandlerFunc
	TagList() gin.HandlerFunc
}
