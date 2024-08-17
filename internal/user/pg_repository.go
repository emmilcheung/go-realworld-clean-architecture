package user

import (
	"context"

	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
)

type Repository interface {
	FindOneUser(ctx context.Context, condition interface{}) (models.User, error)
	SaveOne(c context.Context, data interface{}) error
	Update(c context.Context, data models.User) error
	IsUserFollowing(c context.Context, userId, followerId uint) bool
	GetFollowingsByUser(ctx context.Context, userId uint) []models.User
	SetUserFollow(c context.Context, userId, followerId uint) error
	RemoveUserFollow(c context.Context, userId, followerId uint) error
}
