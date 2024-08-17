package http

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/utils"
)

type ArticleModelValidator struct {
	Article struct {
		Title       string   `form:"title" json:"title" binding:"required,min=4"`
		Description string   `form:"description" json:"description" binding:"max=2048"`
		Body        string   `form:"body" json:"body" binding:"max=2048"`
		Tags        []string `form:"tagList" json:"tagList"`
	} `json:"article"`
}

func NewArticleModelValidator() ArticleModelValidator {
	return ArticleModelValidator{}
}

func (s *ArticleModelValidator) Verify(c *gin.Context) error {
	err := utils.ApplyGinValidator(c, s)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

type ArticlePartialModelValidator struct {
	Article struct {
		Title       string   `form:"title" json:"title" binding:"omitempty,min=4"`
		Description string   `form:"description" json:"description" binding:"omitempty,max=2048"`
		Body        string   `form:"body" json:"body" binding:"omitempty,max=2048"`
		Tags        []string `form:"tagList" json:"omitempty,tagList"`
	} `json:"article"`
}

func (s *ArticlePartialModelValidator) Verify(c *gin.Context) error {
	err := utils.ApplyGinValidator(c, s)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func NewArticlePartialModelValidator() ArticlePartialModelValidator {
	return ArticlePartialModelValidator{}
}

type CommentModelValidator struct {
	Comment struct {
		Body string `form:"body" json:"body" binding:"max=2048"`
	} `json:"comment"`
}

func NewCommentModelValidator() CommentModelValidator {
	return CommentModelValidator{}
}

func (s *CommentModelValidator) Verify(c *gin.Context) error {
	err := utils.ApplyGinValidator(c, s)
	if err != nil {
		return err
	}
	return nil
}
