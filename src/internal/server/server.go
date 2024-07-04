package server

import (
	"fmt"
	"github.com/Eugune-Usachev/social-network/src/internal/handler"
	"github.com/Eugune-Usachev/social-network/src/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HTTPServer struct {
	handler *handler.Handler
	logger  logger.Logger
}

func NewHTTPServer(handler *handler.Handler, logger logger.Logger) *HTTPServer {
	return &HTTPServer{
		handler,
		logger,
	}
}

func (server *HTTPServer) MustStart(addr string, isProduction bool) {
	if isProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	httpServer := http.Server{
		Handler:     server.handler.Handler(),
		Addr:        addr,
		IdleTimeout: time.Minute,
	}

	server.logger.Info(fmt.Sprintf("Starting HTTP server on %s", addr))

	if err := httpServer.ListenAndServe(); err != nil {
		server.logger.Fatal(fmt.Sprintf("Failed to start HTTP server, reason: %s", err))
	}
}
