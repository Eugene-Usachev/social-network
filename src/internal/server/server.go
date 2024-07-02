package server

import (
	"fmt"
	"net/http"
	"social-network/src/internal/handler"
	"social-network/src/pkg/logger"
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

func (server *HTTPServer) MustStart(addr string) {
	if err := http.ListenAndServe(addr, server.handler.Handler()); err != nil {
		server.logger.Fatal(fmt.Sprintf("Failed to start HTTP server, reason: %s", err))
	}
}
