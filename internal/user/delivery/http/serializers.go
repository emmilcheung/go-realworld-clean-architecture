package http

import (
	"context"

	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/user"
)

type ProfileSerializer struct {
	c        context.Context
	userRepo user.Repository
	models.User
}

// Declare your response schema here
type ProfileResponse struct {
	ID        uint    `json:"-"`
	Username  string  `json:"username"`
	Bio       string  `json:"bio"`
	Image     *string `json:"image"`
	Following bool    `json:"following"`
}

// Put your response logic including wrap the userModel here.
func (self *ProfileSerializer) Response() ProfileResponse {
	myUserModel := self.c.Value("my_user_model").(models.User)
	profile := ProfileResponse{
		ID:        self.ID,
		Username:  self.Username,
		Bio:       self.Bio,
		Image:     self.Image,
		Following: self.userRepo.IsUserFollowing(self.c, self.ID, myUserModel.ID),
	}
	return profile
}

type UserSerializer struct {
	c     context.Context
	token string
	models.User
}

type UserResponse struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Bio      string  `json:"bio"`
	Image    *string `json:"image"`
	Token    string  `json:"token"`
}

func (self *UserSerializer) Response() UserResponse {
	user := UserResponse{
		Username: self.Username,
		Email:    self.Email,
		Bio:      self.Bio,
		Image:    self.Image,
		Token:    self.token,
		// Token:    ,
	}
	return user
}
