package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/gothinkster/golang-gin-realworld-example-app/config"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/server"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/db/postgres"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/db/redis"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/locker"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/utils"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func PromHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func Migrate(db *postgres.DB) {
	db.AutoMigrate(&models.Follow{})
	db.AutoMigrate(&models.Article{})
	db.AutoMigrate(&models.Tag{})
	db.AutoMigrate(&models.Favorite{})
	db.AutoMigrate(&models.ArticleUser{})
	db.AutoMigrate(&models.Comment{})
	db.AutoMigrate(&models.User{})
}

func main() {

	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	// postgres
	var sslMode string
	if cfg.Postgres.PostgresqlSSLMode {
		sslMode = "enable"
	} else {
		sslMode = "disable"
	}
	postgresUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Postgres.PostgresqlUser,
		cfg.Postgres.PostgresqlPassword,
		cfg.Postgres.PostgresqlHost,
		cfg.Postgres.PostgresqlPort,
		cfg.Postgres.PostgresqlDbname,
		sslMode,
	)

	gormDB, err := postgres.DBInit(postgresUrl)
	if err != nil {
		log.Fatalf("Postgresql init: %s", err)
	} else {
		log.Printf("Postgres connected, Status: %#v\n", gormDB.DB().Stats())
	}
	defer gormDB.Close()

	Migrate(gormDB)

	// redis
	redis := redis.RedisInit(cfg.Redis.RedisAddr)
	defer redis.Close()

	locker := locker.LockerInit(redis, cfg)

	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: cfg.Jaeger.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           cfg.Jaeger.LogSpans,
			LocalAgentHostPort: cfg.Jaeger.Host,
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatal("cannot create tracer", err)
	}
	log.Println("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	log.Println("Opentracing connected")

	server := server.NewServer(cfg, gormDB, redis, locker)
	if err = server.Run(); err != nil {
		log.Fatal(err)
	}
}
