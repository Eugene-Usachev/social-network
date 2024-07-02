package handler

import (
	"fmt"
	"github.com/Eugene-Usachev/fst"
	"github.com/gin-gonic/gin"
	"net/http"
	"social-network/src/internal/metrics"
	"social-network/src/internal/service"
	"social-network/src/pkg/logger"
	"time"
)

type Handler struct {
	server                *gin.Engine
	logger                logger.Logger
	service               *service.Service
	accessTokenConverter  *fst.EncodedConverter
	refreshTokenConverter *fst.EncodedConverter
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
		"http request | %-7s | %-42s | %d | %-12d microseconds |\n",
		method, path, statusCode, elapsed.Microseconds()))
	metrics.ObserveRequest(elapsed, method, path, statusCode)
}

func (handler *Handler) initRoutes() {
}

func NewHandler(service *service.Service, accessTokenConverter *fst.EncodedConverter, refreshTokenConverter *fst.EncodedConverter) *Handler {
	handler := &Handler{
		service:               service,
		accessTokenConverter:  accessTokenConverter,
		refreshTokenConverter: refreshTokenConverter,
	}

	handler.server = gin.New()
	handler.server.Use(handler.recover, handler.metrics)
	handler.initRoutes()

	return handler
}

func (handler *Handler) Handler() http.Handler {
	return handler.server
}
