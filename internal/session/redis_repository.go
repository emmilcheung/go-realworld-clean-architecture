//go:generate mockgen -source redis_repository.go -destination mock/redis_repository_mock.go -package mock
package session

import (
	"context"

	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
)

// Session repository
type SessRepository interface {
	CreateSession(ctx context.Context, sessionID string, jwt string, expire int) (*models.Session, error)
	GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error)
	DeleteByID(ctx context.Context, sessionID string) error
}
