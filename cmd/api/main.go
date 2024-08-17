package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gothinkster/golang-gin-realworld-example-app/internal/models"
	"github.com/gothinkster/golang-gin-realworld-example-app/internal/server"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/db/postgres"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/db/redis"
	"github.com/gothinkster/golang-gin-realworld-example-app/pkg/locker"

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

	redisUrl := "127.0.0.1:6379"
	postgresUrl := "postgresql://postgres:postgres@localhost:5432/realworld?sslmode=disable"
	gormDB, err := postgres.DBInit(postgresUrl)
	if err != nil {
		log.Fatalf("Postgresql init: %s", err)
	} else {
		log.Printf("Postgres connected, Status: %#v\n", gormDB.DB().Stats())
	}
	defer gormDB.Close()

	Migrate(gormDB)

	redis := redis.RedisInit(redisUrl)
	defer redis.Close()

	locker := locker.LockerInit(redis)

	var ServiceName = "api"
	var LogSpans = false
	var JaegerHost = "localhost:6831"
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           LogSpans,
			LocalAgentHostPort: JaegerHost,
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

	server := server.NewServer(gormDB, redis, locker)
	if err = server.Run(); err != nil {
		log.Fatal(err)
	}
}
