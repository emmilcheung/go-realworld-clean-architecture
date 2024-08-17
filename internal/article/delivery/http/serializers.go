package http

import (
	"context"
	"sort"

	"github.com/gosimple/slug"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/article"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
	userHttp "github.com/gothinkster/golang-gin-realworld-example-app/internal/user/delivery/http"
)

type TagSerializer struct {
	C context.Context
	models.Tag
}

type TagsSerializer struct {
	C    context.Context
	Tags []models.Tag
}

func (s *TagSerializer) Response() string {
	return s.Tag.Tag
}

func (s *TagsSerializer) Response() []string {
	response := []string{}
	for _, tag := range s.Tags {
		serializer := TagSerializer{s.C, tag}
		response = append(response, serializer.Response())
	}
	return response
}

type ArticleUserSerializer struct {
	C context.Context
	models.ArticleUser
}

func (s *ArticleUserSerializer) Response() userHttp.ProfileResponse {
	profile := userHttp.ProfileResponse{
		ID:       s.ArticleUser.User.ID,
		Username: s.ArticleUser.User.Username,
		Bio:      s.ArticleUser.User.Bio,
		Image:    s.ArticleUser.User.Image,
	}
	return profile
}

type ArticleSerializer struct {
	C           context.Context
	articleRepo article.Repository
	models.Article
}

type ArticleResponse struct {
	ID             uint                     `json:"-"`
	Title          string                   `json:"title"`
	Slug           string                   `json:"slug"`
	Description    string                   `json:"description"`
	Body           string                   `json:"body"`
	CreatedAt      string                   `json:"createdAt"`
	UpdatedAt      string                   `json:"updatedAt"`
	Author         userHttp.ProfileResponse `json:"author"`
	Tags           []string                 `json:"tagList"`
	Favorite       bool                     `json:"favorited"`
	FavoritesCount uint                     `json:"favoritesCount"`
}

type ArticlesSerializer struct {
	C           context.Context
	articleRepo article.Repository
	Articles    []models.Article
}

func (s *ArticleSerializer) Response() ArticleResponse {
	myUserModel := s.C.Value("my_user_model").(models.User)
	authorSerializer := ArticleUserSerializer{s.C, s.Author}
	response := ArticleResponse{
		ID:          s.ID,
		Slug:        slug.Make(s.Title),
		Title:       s.Title,
		Description: s.Description,
		Body:        s.Body,
		CreatedAt:   s.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		// UpdatedAt:   s.UpdatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt: s.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		Author:    authorSerializer.Response(),
	}

	articleUserModel := s.articleRepo.GetArticleUser(s.C, myUserModel.ID)
	response.Favorite = s.articleRepo.IsArticleFavoriteBy(s.C, articleUserModel.ID, s.ID)
	response.FavoritesCount = s.articleRepo.ArticleFavoritesCount(s.C, s.ID)

	response.Tags = make([]string, 0)
	sortTags(s.Tags)
	for _, tag := range s.Tags {
		serializer := TagSerializer{s.C, tag}
		response.Tags = append(response.Tags, serializer.Response())
	}
	return response
}

func (s *ArticlesSerializer) Response() []ArticleResponse {
	response := []ArticleResponse{}
	for _, article := range s.Articles {
		serializer := ArticleSerializer{s.C,
			s.articleRepo, article}
		response = append(response, serializer.Response())
	}
	return response
}

func sortTags(tags []models.Tag) {
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Tag < tags[j].Tag
	})
}

type CommentSerializer struct {
	C context.Context
	models.Comment
}

type CommentsSerializer struct {
	C        context.Context
	Comments []models.Comment
}

type CommentResponse struct {
	ID        uint                     `json:"id"`
	Body      string                   `json:"body"`
	CreatedAt string                   `json:"createdAt"`
	UpdatedAt string                   `json:"updatedAt"`
	Author    userHttp.ProfileResponse `json:"author"`
}

func (s *CommentSerializer) Response() CommentResponse {
	response := CommentResponse{
		ID:        s.ID,
		Body:      s.Body,
		CreatedAt: s.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt: s.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
	}
	return response
}

func (s *CommentsSerializer) Response() []CommentResponse {
	response := []CommentResponse{}
	for _, comment := range s.Comments {
		serializer := CommentSerializer{s.C, comment}
		response = append(response, serializer.Response())
	}
	return response
}
