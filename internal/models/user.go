package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint    `gorm:"primary_key"`
	Username     string  `gorm:"column:username"`
	Email        string  `gorm:"column:email;unique_index"`
	Bio          string  `gorm:"column:bio;size:1024"`
	Image        *string `gorm:"column:image"`
	PasswordHash string  `gorm:"column:password;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (e *User) TableName() string {
	return "user_models"
}

// What's bcrypt? https://en.wikipedia.org/wiki/Bcrypt
// Golang bcrypt doc: https://godoc.org/golang.org/x/crypto/bcrypt
// You can change the value in bcrypt.DefaultCost to adjust the security index.
// 	err := userModel.setPassword("password0")
func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.PasswordHash = string(passwordHash)
	return nil
}

// Database will only save the hashed string, you should check it by util function.
// 	if err := serModel.checkPassword("password0"); err != nil { password error }
func (u *User) CheckPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

// A hack way to save ManyToMany relationship,
// gorm will build the alias as FollowingBy <-> FollowingByID <-> "following_by_id".
//
// DB schema looks like: id, created_at, updated_at, deleted_at, following_id, followed_by_id.
//
// Retrieve them by:
// 	db.Where(Follow{ FollowingID:  v.ID, FollowedByID: u.ID, }).First(&follow)
// 	db.Where(Follow{ FollowedByID: u.ID, }).Find(&follows)
//
// More details about gorm.Model: http://jinzhu.me/gorm/models.html#conventions
type Follow struct {
	ID           uint `gorm:"primary_key"`
	Following    User
	FollowingID  uint
	FollowedBy   User
	FollowedByID uint
}

func (e *Follow) TableName() string {
	return "follow_models"
}
