package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gothinkster/golang-gin-realworld-example-app/config"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/session"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/user"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/httpErrors"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/locker"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/utils"
	"github.com/opentracing/opentracing-go"
)

type userHandlers struct {
	cfg      *config.Config
	userRepo user.Repository
	sessUC   session.UseCase
	locker   locker.Locker
}

func NewUserHandlers(cfg *config.Config, userRepo user.Repository, sessUC session.UseCase, locker locker.Locker) user.Handlers {
	return &userHandlers{cfg, userRepo, sessUC, locker}
}

func (h userHandlers) UsersRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "user.UsersRegistration")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		userModelValidator := NewUserModelValidator()
		if err := userModelValidator.Bind(c); err != nil {
			fmt.Print(err)
			c.JSON(http.StatusUnprocessableEntity, httpErrors.NewValidatorError(err))
			return
		}
		lockKey := fmt.Sprintf("user:email-%s", userModelValidator.userModel.Email)
		lock := h.locker.ObtainLock(ctx, lockKey)
		defer lock.Release(ctx)

		if err := h.userRepo.SaveOne(ctx, &userModelValidator.userModel); err != nil {
			c.JSON(http.StatusUnprocessableEntity, httpErrors.NewError("database", err))
			return
		}
		userSession, _ := h.sessUC.CreateSession(ctx, &userModelValidator.userModel, h.cfg.Session.Expire)
		serializer := UserSerializer{ctx, userSession.Token, userModelValidator.userModel}

		c.JSON(http.StatusCreated, gin.H{"user": serializer.Response()})
	}
}

func (h userHandlers) UsersLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "user.UsersLogin")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		loginValidator := NewLoginValidator()
		if err := loginValidator.Bind(c); err != nil {
			c.JSON(http.StatusUnprocessableEntity, httpErrors.NewValidatorError(err))
			return
		}
		userModel, err := h.userRepo.FindOneUser(ctx, &models.User{Email: loginValidator.userModel.Email})

		if err != nil {
			c.JSON(http.StatusForbidden, httpErrors.NewError("login", errors.New("Not Registered email or invalid password")))
			return
		}

		if userModel.CheckPassword(loginValidator.User.Password) != nil {
			c.JSON(http.StatusForbidden, httpErrors.NewError("login", errors.New("Not Registered email or invalid password")))
			return
		}
		userSession, _ := h.sessUC.CreateSession(ctx, &userModel, h.cfg.Session.Expire)
		serializer := UserSerializer{ctx, userSession.Token, userModel}
		c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
	}
}

func (h userHandlers) UserRetrieve() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "user.UserRetrieve")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		userModel := c.MustGet("my_user_model").(models.User)
		sessionID := c.MustGet("my_session_id").(string)
		userSession, _ := h.sessUC.GetSessionByID(ctx, sessionID)
		serializer := UserSerializer{ctx, userSession.Token, userModel}
		c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
	}
}

func (h userHandlers) UserUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "user.UserUpdate")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		myUserModel := c.MustGet("my_user_model").(models.User)
		userModelValidator := NewUserModelValidatorFillWith(myUserModel)
		if err := userModelValidator.Bind(c); err != nil {
			c.JSON(http.StatusUnprocessableEntity, httpErrors.NewValidatorError(err))
			return
		}

		lockKey := fmt.Sprintf("user:email-%s", userModelValidator.userModel.Email)
		lock := h.locker.ObtainLock(ctx, lockKey)
		defer lock.Release(ctx)

		userModelValidator.userModel.ID = myUserModel.ID
		if err := h.userRepo.Update(ctx, userModelValidator.userModel); err != nil {
			c.JSON(http.StatusUnprocessableEntity, httpErrors.NewError("database", err))
			return
		}
		userModel := c.MustGet("my_user_model").(models.User)
		sessionID := c.MustGet("my_session_id").(string)
		userSession, _ := h.sessUC.GetSessionByID(ctx, sessionID)
		serializer := UserSerializer{ctx, userSession.Token, userModel}

		c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
	}
}

func (h userHandlers) ProfileRetrieve() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "user.ProfileRetrieve")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		username := c.Param("username")
		userModel, err := h.userRepo.FindOneUser(ctx, &models.User{Username: username})
		if err != nil {
			c.JSON(http.StatusNotFound, httpErrors.NewError("profile", errors.New("Invalid username")))
			return
		}
		profileSerializer := ProfileSerializer{ctx, h.userRepo, userModel}
		c.JSON(http.StatusOK, gin.H{"profile": profileSerializer.Response()})
	}
}
func (h userHandlers) ProfileFollow() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "user.ProfileFollow")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		username := c.Param("username")
		userModel, err := h.userRepo.FindOneUser(ctx, &models.User{Username: username})
		if err != nil {
			c.JSON(http.StatusNotFound, httpErrors.NewError("profile", errors.New("Invalid username")))
			return
		}
		myUserModel := c.MustGet("my_user_model").(models.User)
		err = h.userRepo.SetUserFollow(ctx, userModel.ID, myUserModel.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, httpErrors.NewError("database", err))
			return
		}
		serializer := ProfileSerializer{ctx, h.userRepo, userModel}
		c.JSON(http.StatusOK, gin.H{"profile": serializer.Response()})
	}
}
func (h userHandlers) ProfileUnfollow() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "user.ProfileUnfollow")
		span.SetTag("requestId", utils.GetRequestID(c))
		defer span.Finish()

		username := c.Param("username")
		userModel, err := h.userRepo.FindOneUser(ctx, &models.User{Username: username})
		if err != nil {
			c.JSON(http.StatusNotFound, httpErrors.NewError("profile", errors.New("Invalid username")))
			return
		}
		myUserModel := c.MustGet("my_user_model").(models.User)

		err = h.userRepo.RemoveUserFollow(ctx, userModel.ID, myUserModel.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, httpErrors.NewError("database", err))
			return
		}
		serializer := ProfileSerializer{ctx, h.userRepo, userModel}
		c.JSON(http.StatusOK, gin.H{"profile": serializer.Response()})
	}
}
