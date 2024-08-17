package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/session"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/utils"
	"github.com/opentracing/opentracing-go"
)

// Session use case
type sessionUC struct {
	sessionRepo session.SessRepository
}

// New session use case constructor
func NewSessionUseCase(sessionRepo session.SessRepository) session.UseCase {
	return &sessionUC{sessionRepo: sessionRepo}
}

// Create new session
func (u *sessionUC) CreateSession(ctx context.Context, user *models.User, expire int) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.sessionUC.CreateSession")
	defer span.Finish()

	sessionID := uuid.New().String()
	jwt := utils.GenToken(user.ID, sessionID)
	return u.sessionRepo.CreateSession(ctx, sessionID, jwt, expire)
}

// Delete session by id
func (u *sessionUC) DeleteByID(ctx context.Context, sessionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.sessionUC.DeleteByID")
	defer span.Finish()

	return u.sessionRepo.DeleteByID(ctx, sessionID)
}

// get session by id
func (u *sessionUC) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.sessionUC.GetSessionByID")
	defer span.Finish()

	return u.sessionRepo.GetSessionByID(ctx, sessionID)
}
