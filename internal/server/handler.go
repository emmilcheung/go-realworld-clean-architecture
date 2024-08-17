package server

import (
	"github.com/gin-gonic/gin"
	requestid "github.com/sumit-tembe/gin-requestid"

	articleHttp "github.com/gothinkster/golang-gin-realworld-example-app/internal/article/delivery/http"
	articleRepository "github.com/gothinkster/golang-gin-realworld-example-app/internal/article/repository"
	articleUsecase "github.com/gothinkster/golang-gin-realworld-example-app/internal/article/usecase"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/middleware"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/metric"

	sessRepository "github.com/gothinkster/golang-gin-realworld-example-app/internal/session/repository"

	sessUsecase "github.com/gothinkster/golang-gin-realworld-example-app/internal/session/usecase"

	userHttp "github.com/gothinkster/golang-gin-realworld-example-app/internal/user/delivery/http"

	userRepository "github.com/gothinkster/golang-gin-realworld-example-app/internal/user/repository"
)

func (s *Server) MapHandlers(engine *gin.Engine) error {

	// resources
	articleRepo := articleRepository.NewArticleRepository(s.db)
	articleUc := articleUsecase.NewArticleUseCase(articleRepo)
	articleHandlers := articleHttp.NewArticleHandlers(articleRepo, articleUc, s.locker)
	sessRepo := sessRepository.NewSessionRepository(s.redisClient)
	sessUC := sessUsecase.NewSessionUseCase(sessRepo)
	userRepo := userRepository.NewUserRepository(s.db)
	userHandler := userHttp.NewUserHandlers(userRepo, sessUC, s.locker)
	mv := middleware.NewMiddlewareManager(s.db, sessUC, []string{"*"})

	// Middlewares
	{
		p := metric.NewPrometheus("gin")
		p.Use(engine)
		engine.Use(mv.MetricsMiddleware(p))
		//recovery middleware
		engine.Use(gin.Recovery())
		//middleware which injects a 'RequestID' into the context and header of each request.
		engine.Use(requestid.RequestID(nil))
		//middleware which enhance Gin request logger to include 'RequestID'
		engine.Use(gin.LoggerWithConfig(requestid.GetLoggerConfig(nil, nil, nil)))
	}

	// routes
	{
		v1 := engine.Group("/api")
		// public routes
		v1.Use(mv.AuthMiddleware(s.db, false))
		userHttp.UsersRegister(v1.Group("/users"), userHandler)
		articleHttp.ArticlesAnonymousRouteRegister(v1.Group("/articles"), articleHandlers)
		articleHttp.TagsAnonymousRouteRegister(v1.Group("/tags"), articleHandlers)

		// protected routes
		v1.Use(mv.AuthMiddleware(s.db, true))
		articleHttp.ArticlesRouteRegister(v1.Group("/articles"), articleHandlers)
		userHttp.UserRegister(v1.Group("/user"), userHandler)
		userHttp.ProfileRegister(v1.Group("/profiles"), userHandler)
	}

	return nil
}
