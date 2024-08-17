package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/user"
)

func UsersRegister(router *gin.RouterGroup, h user.Handlers) {
	router.POST("/", h.UsersRegistration())
	router.POST("/login", h.UsersLogin())
}

func UserRegister(router *gin.RouterGroup, h user.Handlers) {
	router.GET("/", h.UserRetrieve())
	router.PUT("/", h.UserUpdate())
}

func ProfileRegister(router *gin.RouterGroup, h user.Handlers) {
	router.GET("/:username", h.ProfileRetrieve())
	router.POST("/:username/follow", h.ProfileFollow())
	router.DELETE("/:username/follow", h.ProfileUnfollow())
}
