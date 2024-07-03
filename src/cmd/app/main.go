package main

import (
	"fmt"
	fb "github.com/Eugene-Usachev/fastbytes"
	"github.com/Eugene-Usachev/fst"
	"social-network/src/internal/config"
	"social-network/src/internal/handler"
	"social-network/src/internal/metrics"
	"social-network/src/internal/repository"
	"social-network/src/internal/server"
	"social-network/src/internal/service"
	loggerpkg "social-network/src/pkg/logger"
	"time"
)

func main() {
	cfg := config.MustNewConfig()

	logger := loggerpkg.MustCreateLogger(cfg.IsProduction(), cfg)

	logger.Info("Created logger")

	accessTokenConverter := fst.NewEncodedConverter(&fst.ConverterConfig{
		ExpirationTime: time.Second * 300,
		SecretKey:      fb.S2B(cfg.FstAccessKey()),
	})
	refreshTokenConverter := fst.NewEncodedConverter(&fst.ConverterConfig{
		ExpirationTime: time.Hour * 24 * 30,
		SecretKey:      fb.S2B(cfg.FstRefreshKey()),
	})

	go metrics.MustListen(cfg.MetricsAddr(), logger)

	repository := repository.NewRepository()
	service := service.NewService()
	handler := handler.NewHandler(service, accessTokenConverter, refreshTokenConverter)
	server := server.NewHTTPServer(handler, logger)

	server.MustStart(fmt.Sprintf("%s:%d", cfg.Host(), cfg.Port()))
}
