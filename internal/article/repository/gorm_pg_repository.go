package repository

import (
	"context"

	article "github.com/gothinkster/golang-gin-realworld-example-app/internal/article"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
	UserRepo "github.com/gothinkster/golang-gin-realworld-example-app/internal/user/repository"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/db/postgres"
	"github.com/opentracing/opentracing-go"
)

type articleRepo struct {
	article.Repository
	db *postgres.DB
}

func NewArticleRepository(db *postgres.DB) article.Repository {
	return &articleRepo{db: db}
}

func (r *articleRepo) GetArticleUser(c context.Context, userID uint) models.ArticleUser {
	span, _ := opentracing.StartSpanFromContext(c, "article.articleRepo.GetArticleUser")
	defer span.Finish()
	var articleUserModel models.ArticleUser
	var userModel models.User
	if userID == 0 {
		return articleUserModel
	}
	db := r.db
	db.Where(&models.User{ID: userID}).First(&userModel)
	db.Where(&models.ArticleUser{
		UserID: userID,
	}).FirstOrCreate(&articleUserModel)

	articleUserModel.User = userModel
	return articleUserModel
}

func (r *articleRepo) FindOneArticle(c context.Context, condition interface{}) (models.Article, error) {
	span, _ := opentracing.StartSpanFromContext(c, "article.articleRepo.FindOneArticle")
	defer span.Finish()

	db := r.db
	var model models.Article
	tx := db.Begin()
	tx.Where(condition).First(&model)
	tx.Model(&model).Related(&model.Author, "Author")
	tx.Model(&model.Author).Related(&model.Author.User, "UserModelID")
	tx.Model(&model).Related(&model.Tags, "Tags")
	err := tx.Commit().Error
	return model, err
}

func (r *articleRepo) ArticleFavoritesCount(c context.Context, articleId uint) uint {
	span, _ := opentracing.StartSpanFromContext(c, "article.articleRepo.ArticleFavoritesCount")
	defer span.Finish()

	db := r.db
	var count uint
	db.Model(&models.Favorite{}).Where(models.Favorite{
		FavoriteID: articleId,
	}).Count(&count)
	return count
}

func (r *articleRepo) IsArticleFavoriteBy(c context.Context, userId uint, articleId uint) bool {
	span, _ := opentracing.StartSpanFromContext(c, "article.articleRepo.IsArticleFavoriteBy")
	defer span.Finish()

	db := r.db
	var favorite models.Favorite
	db.Where(models.Favorite{
		FavoriteID:   articleId,
		FavoriteByID: userId,
	}).First(&favorite)
	return favorite.ID != 0
}

func (r *articleRepo) FindManyArticle(ctx context.Context, tag, author, favorited string, limit, offset int) ([]models.Article, int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "article.articleRepo.FindManyArticle")
	defer span.Finish()
	db := r.db
	var articleModels []models.Article
	var count int

	tx := db.Begin()
	if tag != "" {
		var tagModel models.Tag
		tx.Where(models.Tag{Tag: tag}).First(&tagModel)
		if tagModel.ID != 0 {
			tx.Model(&tagModel).Offset(offset).Limit(limit).Related(&articleModels, "Articles")
			count = tx.Model(&tagModel).Association("Articles").Count()
		}
	} else if author != "" {
		var userModel models.User
		tx.Where(models.User{Username: author}).First(&userModel)
		articleUserModel := r.GetArticleUser(ctx, userModel.ID)

		if articleUserModel.ID != 0 {
			count = tx.Model(&articleUserModel).Association("Articles").Count()
			tx.Model(&articleUserModel).Offset(offset).Limit(limit).Related(&articleModels, "Articles")
		}
	} else if favorited != "" {
		var userModel models.User
		tx.Where(models.User{Username: favorited}).First(&userModel)
		articleUserModel := r.GetArticleUser(ctx, userModel.ID)
		if articleUserModel.ID != 0 {
			var favoriteModels []models.Favorite
			tx.Where(models.Favorite{
				FavoriteByID: articleUserModel.ID,
			}).Offset(offset).Limit(limit).Find(&favoriteModels)

			count = tx.Model(&articleUserModel).Association("Favorites").Count()
			for _, favorite := range favoriteModels {
				var model models.Article
				tx.Model(&favorite).Related(&model, "Favorite")
				articleModels = append(articleModels, model)
			}
		}
	} else {
		db.Model(&articleModels).Count(&count)
		db.Offset(offset).Limit(limit).Find(&articleModels)
	}

	for i, _ := range articleModels {
		tx.Model(&articleModels[i]).Related(&articleModels[i].Author, "Author")
		tx.Model(&articleModels[i].Author).Related(&articleModels[i].Author.User)
		tx.Model(&articleModels[i]).Related(&articleModels[i].Tags, "Tags")
	}
	err := tx.Commit().Error
	return articleModels, count, err
}

func (r *articleRepo) GetArticleFeed(c context.Context, userId uint, limit, offset int) ([]models.Article, int, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "article.articleRepo.GetArticleFeed")
	defer span.Finish()
	var articleModels []models.Article
	var count int

	tx := r.db.Begin()
	userRepo := UserRepo.NewUserRepository(r.db)
	followings := userRepo.GetFollowingsByUser(ctx, userId)
	var articleUserModels []uint
	for _, following := range followings {
		articleUserModel := r.GetArticleUser(ctx, following.ID)
		articleUserModels = append(articleUserModels, articleUserModel.ID)
	}

	tx.Where("author_id in (?)", articleUserModels).Order("updated_at desc").Offset(offset).Limit(limit).Find(&articleModels)

	for i, _ := range articleModels {
		tx.Model(&articleModels[i]).Related(&articleModels[i].Author, "Author")
		tx.Model(&articleModels[i].Author).Related(&articleModels[i].Author.User)
		tx.Model(&articleModels[i]).Related(&articleModels[i].Tags, "Tags")
	}
	err := tx.Commit().Error
	return articleModels, count, err
}

func (r *articleRepo) SaveOne(ctx context.Context, data interface{}) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "article.articleRepo.SaveOne")
	defer span.Finish()

	err := r.db.Save(data).Error
	return err
}

func (r *articleRepo) Update(c context.Context, data *models.Article) error {
	span, _ := opentracing.StartSpanFromContext(c, "article.articleRepo.Update")
	defer span.Finish()
	err := r.db.Model(&models.Article{ID: data.ID}).Update(data).Error
	return err
}

func (r *articleRepo) DeleteArticleModel(c context.Context, condition interface{}) error {
	span, _ := opentracing.StartSpanFromContext(c, "article.articleRepo.DeleteArticleModel")
	defer span.Finish()
	err := r.db.Unscoped().Where(condition).Delete(models.Article{}).Error
	return err
}

func (r *articleRepo) UpsertTags(c context.Context, tags []string) ([]models.Tag, error) {
	span, _ := opentracing.StartSpanFromContext(c, "article.articleRepo.UpsertTags")
	defer span.Finish()

	var tagList []models.Tag
	for _, tag := range tags {
		var tagModel models.Tag
		err := r.db.FirstOrCreate(&tagModel, models.Tag{Tag: tag}).Error
		if err != nil {
			return nil, err
		}
		tagList = append(tagList, tagModel)
	}
	return tagList, nil
}

func (r *articleRepo) SetFavorite(c context.Context, articleId, userId uint) error {
	span, _ := opentracing.StartSpanFromContext(c, "article.articleRepo.SetFavorite")
	defer span.Finish()

	var favorite models.Favorite
	err := r.db.FirstOrCreate(&favorite, &models.Favorite{
		FavoriteID:   articleId,
		FavoriteByID: userId,
	}).Error
	return err
}

func (r *articleRepo) RemoveFavorite(c context.Context, articleId, userId uint) error {
	span, _ := opentracing.StartSpanFromContext(c, "article.articleRepo.RemoveFavorite")
	defer span.Finish()

	err := r.db.Unscoped().Where(models.Favorite{
		FavoriteID:   articleId,
		FavoriteByID: userId,
	}).Delete(&models.Favorite{}).Error
	return err
}

func (r *articleRepo) GetArticleComments(c context.Context, article models.Article) ([]models.Comment, error) {
	span, _ := opentracing.StartSpanFromContext(c, "article.articleRepo.GetArticleComments")
	defer span.Finish()

	tx := r.db.Begin()
	tx.Model(&article).Related(&article.Comments, "Comments")
	for i, _ := range article.Comments {
		tx.Model(&article.Comments[i]).Related(&article.Comments[i].Author, "Author")
		tx.Model(&article.Comments[i].Author).Related(&article.Comments[i].Author.User)
	}
	err := tx.Commit().Error
	return article.Comments, err
}

func (r *articleRepo) DeleteComment(c context.Context, condition interface{}) error {
	span, _ := opentracing.StartSpanFromContext(c, "article.articleRepo.DeleteComment")
	defer span.Finish()

	err := r.db.Unscoped().Where(condition).Delete(models.Comment{}).Error
	return err
}

func (r *articleRepo) GetTags() ([]models.Tag, error) {
	var models []models.Tag
	err := r.db.Find(&models).Error
	return models, err
}
