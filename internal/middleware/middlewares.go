package middleware

import (
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/session"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/db/postgres"
)

// Middleware manager
type MiddlewareManager struct {
	db     *postgres.DB
	sessUC session.UseCase
	// authUC auth.UseCase
	// cfg     *config.Config
	origins []string
	// logger  logger.Logger
}

// Middleware manager constructor
func NewMiddlewareManager(db *postgres.DB, sessUC session.UseCase, origins []string) *MiddlewareManager {
	return &MiddlewareManager{db: db, sessUC: sessUC, origins: origins}
}
