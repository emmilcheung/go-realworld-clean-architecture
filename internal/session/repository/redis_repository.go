package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gothinkster/golang-gin-realworld-example-app/config"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/session"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/db/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// Session repository
type sessionRepo struct {
	redisClient *redis.Client
	appName     string
	basePrefix  string
}

// Session repository constructor
func NewSessionRepository(cfg *config.Config, redisClient *redis.Client) session.SessRepository {
	return &sessionRepo{redisClient: redisClient, appName: cfg.Server.AppName, basePrefix: fmt.Sprintf("%s:", cfg.Session.Prefix)}
}

// Create session in redis
func (s *sessionRepo) CreateSession(ctx context.Context, sessionID string, jwt string, expire int) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.sessionRepo.CreateSession")
	defer span.Finish()

	var sess models.Session
	sess.SessionID = sessionID
	sess.Token = jwt
	sessionKey := s.buildKey(sess.SessionID)

	sessBytes, err := json.Marshal(&sess)
	if err != nil {
		return nil, errors.WithMessage(err, "sessionRepo.CreateSession.json.Marshal")
	}
	if err = s.redisClient.Set(ctx, sessionKey, sessBytes, time.Second*time.Duration(expire)).Err(); err != nil {
		return nil, errors.Wrap(err, "sessionRepo.CreateSession.redisClient.Set")
	}
	return &sess, nil
}

// Get session by id
func (s *sessionRepo) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.sessionRepo.GetSessionByID")
	defer span.Finish()

	sessBytes, err := s.redisClient.Get(ctx, s.buildKey(sessionID)).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "sessionRep.GetSessionByID.redisClient.Get")
	}

	sess := &models.Session{}
	if err = json.Unmarshal(sessBytes, &sess); err != nil {
		return nil, errors.Wrap(err, "sessionRepo.GetSessionByID.json.Unmarshal")
	}
	return sess, nil
}

// Delete session by id
func (s *sessionRepo) DeleteByID(ctx context.Context, sessionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "session.sessionRepo.DeleteByID")
	defer span.Finish()

	if err := s.redisClient.Del(ctx, s.buildKey(sessionID)).Err(); err != nil {
		return errors.Wrap(err, "sessionRepo.DeleteByID")
	}
	return nil
}

func (s *sessionRepo) buildKey(sessionID string) string {
	return fmt.Sprintf("%s:%s:%s", s.appName, s.basePrefix, sessionID)
}
