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
	"github.com/gothinkster/golang-gin-realworld-example-app/config"
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
	cfg         *config.Config
	// logger      logger.Logger
}

// NewServer New Server constructor
func NewServer(cfg *config.Config, db *postgres.DB, redisClient *redis.Client, locker locker.Locker) *Server {
	var serverMode string
	if cfg.Server.Debug {
		serverMode = gin.DebugMode
	} else {
		serverMode = gin.ReleaseMode
	}
	gin.SetMode(serverMode)
	return &Server{gin: gin.Default(), cfg: cfg, db: db, redisClient: redisClient, locker: locker}
}

func (s *Server) Run() error {
	if err := s.MapHandlers(s.gin); err != nil {
		return err
	}

	// Ref: https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
	srv := &http.Server{
		Addr:           s.cfg.Server.Port,
		Handler:        s.gin.Handler(),
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
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
