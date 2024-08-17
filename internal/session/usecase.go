//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package session

import (
	"context"

	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
)

// Session use case
type UseCase interface {
	CreateSession(ctx context.Context, user *models.User, expire int) (*models.Session, error)
	GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error)
	DeleteByID(ctx context.Context, sessionID string) error
}
