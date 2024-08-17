package user

import "github.com/gin-gonic/gin"

type Handlers interface {
	UsersRegistration() gin.HandlerFunc
	UsersLogin() gin.HandlerFunc
	UserRetrieve() gin.HandlerFunc
	UserUpdate() gin.HandlerFunc
	ProfileRetrieve() gin.HandlerFunc
	ProfileFollow() gin.HandlerFunc
	ProfileUnfollow() gin.HandlerFunc
}
