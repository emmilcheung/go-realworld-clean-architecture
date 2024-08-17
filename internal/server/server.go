package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/db/postgres"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/db/redis"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/locker"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	gin         *gin.Engine
	db          *postgres.DB
	redisClient *redis.Client
	locker      locker.Locker
	// cfg         *config.Config
	// logger      logger.Logger
}

// NewServer New Server constructor
func NewServer(db *postgres.DB, redisClient *redis.Client, locker locker.Locker) *Server {
	return &Server{gin: gin.Default(), db: db, redisClient: redisClient, locker: locker}
}

func (s *Server) Run() error {
	if err := s.MapHandlers(s.gin); err != nil {
		return err
	}

	// Ref: https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.gin.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()
	return srv.Shutdown(ctx)
}
