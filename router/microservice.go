package router

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

const (
	dbName         = "users"
	collectionName = "users"
	serviceName    = "users-service"
)

type IMicroservice interface {
	StartHTTP()

	// HTTP Services
	USE(handler ServiceHandleFunc)
	GET(path string, h ServiceHandleFunc)
	POST(path string, h ServiceHandleFunc)
	PUT(path string, h ServiceHandleFunc)
	PATCH(path string, h ServiceHandleFunc)
	DELETE(path string, h ServiceHandleFunc)
}

type Microservice struct {
	*gin.Engine
}

type ServiceHandleFunc func(c IContext)

func NewMicroservice() IMicroservice {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(LoggingMiddleware())

	return &Microservice{r}
}

func (ms *Microservice) USE(handler ServiceHandleFunc) {
	ms.Engine.Use(func(ctx *gin.Context) {
		handler(NewContext(ms, ctx))
	})
}

func (ms *Microservice) GET(path string, handler ServiceHandleFunc) {
	ms.Engine.GET(path, func(ctx *gin.Context) {
		handler(NewContext(ms, ctx))
	})
}

func (ms *Microservice) POST(path string, handler ServiceHandleFunc) {
	ms.Engine.POST(path, func(ctx *gin.Context) {
		handler(NewContext(ms, ctx))
	})
}

func (ms *Microservice) PUT(path string, h ServiceHandleFunc) {
	ms.Engine.PUT(path, func(ctx *gin.Context) {
		h(NewContext(ms, ctx))
	})
}

func (ms *Microservice) PATCH(path string, h ServiceHandleFunc) {
	ms.Engine.PATCH(path, func(ctx *gin.Context) {
		h(NewContext(ms, ctx))
	})
}

func (ms *Microservice) DELETE(path string, handler ServiceHandleFunc) {
	ms.Engine.DELETE(path, func(ctx *gin.Context) {
		handler(NewContext(ms, ctx))
	})
}

func (ms *Microservice) StartHTTP() {
	addr := os.Getenv("PORT")
	if addr == "" {
		log.Error("PORT not found")
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:           ":" + addr,
		Handler:        ms.Engine,
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
