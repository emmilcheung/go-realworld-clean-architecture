package repository

import (
	"context"

	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/user"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/db/postgres"
	"github.com/opentracing/opentracing-go"
)

type userRepo struct {
	user.Repository
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) user.Repository {
	return &userRepo{db: db}
}

func (r *userRepo) FindOneUser(c context.Context, condition interface{}) (models.User, error) {
	span, _ := opentracing.StartSpanFromContext(c, "user.userRepo.FindOneUser")
	defer span.Finish()

	var model models.User
	err := r.db.Where(condition).First(&model).Error
	return model, err
}

// You could input an UserModel which will be saved in database returning with error info
// 	if err := SaveOne(&userModel); err != nil { ... }
func (r *userRepo) SaveOne(c context.Context, data interface{}) error {
	span, _ := opentracing.StartSpanFromContext(c, "user.userRepo.SaveOne")
	defer span.Finish()

	err := r.db.Save(data).Error
	return err
}

func (r *userRepo) Update(c context.Context, data models.User) error {
	span, _ := opentracing.StartSpanFromContext(c, "user.userRepo.Update")
	defer span.Finish()

	err := r.db.Model(&models.User{ID: data.ID}).Update(data).Error
	return err
}

// You could get a following list of userModel
// 	followings := userModel.GetFollowings()
func (r *userRepo) GetFollowingsByUser(c context.Context, userId uint) []models.User {
	span, _ := opentracing.StartSpanFromContext(c, "user.userRepo.GetFollowingsByUser")
	defer span.Finish()

	var tx *postgres.DB
	tx = r.db.Begin()

	var follows []models.Follow
	var followings []models.User
	tx.Where(models.Follow{
		FollowedByID: userId,
	}).Find(&follows)
	for _, follow := range follows {
		var userModel models.User
		tx.Model(&follow).Related(&userModel, "Following")
		followings = append(followings, userModel)
	}
	defer tx.Commit()

	return followings
}

// You could check whether  userModel1 following userModel2
// 	followingBool = myUserModel.isFollowing(self.UserModel)
func (r *userRepo) IsUserFollowing(c context.Context, userId, followerId uint) bool {
	span, _ := opentracing.StartSpanFromContext(c, "user.userRepo.IsUserFollowing")
	defer span.Finish()

	var follow models.Follow
	r.db.Where(models.Follow{
		FollowingID:  userId,
		FollowedByID: followerId,
	}).First(&follow)
	return follow.ID != 0
}

// You could add a following relationship as userModel1 following userModel2
// 	err = userModel1.following(userModel2)
func (r *userRepo) SetUserFollow(c context.Context, userId, followerId uint) error {
	span, _ := opentracing.StartSpanFromContext(c, "user.userRepo.SetUserFollow")
	defer span.Finish()

	var follow models.Follow
	err := r.db.FirstOrCreate(&follow, &models.Follow{
		FollowingID:  userId,
		FollowedByID: followerId,
	}).Error
	return err
}

// You could delete a following relationship as userModel1 following userModel2
// 	err = userModel1.unFollowing(userModel2)
func (r *userRepo) RemoveUserFollow(c context.Context, userId, followerId uint) error {
	span, _ := opentracing.StartSpanFromContext(c, "user.userRepo.RemoveUserFollow")
	defer span.Finish()

	err := r.db.Unscoped().Where(models.Follow{
		FollowingID:  userId,
		FollowedByID: followerId,
	}).Delete(models.Follow{}).Error
	return err
}
