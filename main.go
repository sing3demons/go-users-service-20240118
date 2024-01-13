package main

import (
	"context"
	"errors"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sing3demons/users/handler"
	"github.com/sing3demons/users/middleware"
	"github.com/sing3demons/users/repository"
	"github.com/sing3demons/users/service"
	log "github.com/sirupsen/logrus"
)

func init() {
	if os.Getenv("ZONE") == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		if err := godotenv.Load(".env.dev"); err != nil {
			panic(err)
		}
	}
}

const (
	dbName         = "users"
	collectionName = "users"
	serviceName    = "users-service"
)

func main() {
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: time.RFC3339, PrettyPrint: true})

	db := NewDatabase(dbName, collectionName)
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())

	r.GET("/healthz", healthz)

	r.POST("/auth/register", userHandler.Register)
	r.POST("/auth/login", userHandler.Login)

	r.GET("/profile", middleware.Authorization(), userHandler.GetProfile)

	runServer(r)

}

func healthz(c *gin.Context) {
	c.Status(http.StatusOK)
}

func runServer(router http.Handler) {
	addr := os.Getenv("PORT")
	if addr == "" {
		log.Error("PORT not found")
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:           ":" + addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	go func() {
		log.WithFields(log.Fields{
			"PORT":        srv.Addr,
			"TYPE":        "HTTP",
			"func":        "runServer",
			"APP_NAME":    serviceName,
			"APP_VERSION": "1.0.0",
		}).Info("HTTP server is running")

		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Errorf("server listen err: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("server forced to shutdown: %v", err)
		os.Exit(1)
	}

	log.Info("server exited")
}
