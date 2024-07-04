package handler

import (
	"fmt"
	"github.com/Eugene-Usachev/fst"
	"github.com/gin-gonic/gin"
	"net/http"
	"social-network/src/internal/metrics"
	"social-network/src/internal/service"
	loggerpkg "social-network/src/pkg/logger"
	"time"
)

type Handler struct {
	router                *gin.Engine
	logger                loggerpkg.Logger
	service               *service.Service
	accessTokenConverter  *fst.EncodedConverter
	refreshTokenConverter *fst.EncodedConverter
}

func NewHandler(
	service *service.Service,
	accessTokenConverter *fst.EncodedConverter,
	refreshTokenConverter *fst.EncodedConverter,
	logger loggerpkg.Logger,
) *Handler {
	handler := &Handler{
		service:               service,
		logger:                logger,
		accessTokenConverter:  accessTokenConverter,
		refreshTokenConverter: refreshTokenConverter,
	}

	handler.router = gin.New()
	handler.router.Use(handler.recover, handler.metrics)
	handler.initRoutes()

	return handler
}

func (handler *Handler) recover(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			handler.logger.Error(fmt.Sprintf("Panic recovered in HTTP handler, reason: %v", err))
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
	ctx.Next()
}

func (handler *Handler) metrics(ctx *gin.Context) {
	startTime := time.Now()
	ctx.Next()
	elapsed := time.Since(startTime)
	method := ctx.Request.Method
	path := ctx.Request.URL.Path
	statusCode := ctx.Writer.Status()
	handler.logger.Info(fmt.Sprintf(
		"http request | %-7s | %-21s | %d | %-9d microseconds",
		method, path, statusCode, elapsed.Microseconds()))
	metrics.ObserveRequest(elapsed, method, path, statusCode)
}

func (handler *Handler) initRoutes() {
	authGroup := handler.router.Group("/auth")
	{
		authGroup.POST("/sign-up", handler.SingUp)
		authGroup.POST("/sign-in", handler.SignIn)
		authGroup.POST("/refresh-tokens", handler.Refresh)
	}

	metricsHandler := metrics.Handler()
	handler.router.GET("/metrics", func(ctx *gin.Context) {
		metricsHandler.ServeHTTP(ctx.Writer, ctx.Request)
	})
}

func (handler *Handler) Handler() http.Handler {
	return handler.router
}

func (handler *Handler) AbortWithError(ctx *gin.Context, status int, err error) {
	ctx.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
}
