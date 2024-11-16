package main

import (
	"context"
	"fmt"
	fb "github.com/Eugene-Usachev/fastbytes"
	"github.com/Eugene-Usachev/fst"
	"github.com/Eugune-Usachev/social-network/src/internal/config"
	handlerpkg "github.com/Eugune-Usachev/social-network/src/internal/handler"
	repositorypkg "github.com/Eugune-Usachev/social-network/src/internal/repository"
	"github.com/Eugune-Usachev/social-network/src/internal/repository/cache"
	"github.com/Eugune-Usachev/social-network/src/internal/repository/postgres"
	serverpkg "github.com/Eugune-Usachev/social-network/src/internal/server"
	servicepkg "github.com/Eugune-Usachev/social-network/src/internal/service"
	loggerpkg "github.com/Eugune-Usachev/social-network/src/pkg/logger"
	"strconv"
	"time"
)

func main() {
	cfg := config.MustNewConfig()

	logger := loggerpkg.MustCreateLogger(cfg.IsProduction(), cfg)

	accessTokenConverter := fst.NewEncodedConverter(&fst.ConverterConfig{
		ExpirationTime: time.Second * 300,
		SecretKey:      fb.S2B(cfg.FstAccessKey()),
	})
	refreshTokenConverter := fst.NewEncodedConverter(&fst.ConverterConfig{
		ExpirationTime: time.Hour * 24 * 30,
		SecretKey:      fb.S2B(cfg.FstRefreshKey()),
	})

	repository := repositorypkg.NewRepository(postgres.MustCreatePostgresDB(context.Background(), postgres.Config{
		Host:     cfg.PostgresHost(),
		Port:     strconv.Itoa(cfg.PostgresPort()),
		UserName: cfg.PostgresUser(),
		UserPass: cfg.PostgresPass(),
		DBName:   cfg.PostgresDBName(),
		SSLMode:  cfg.PostgresSSLMode(),
	}, logger), cache.MustCreateRedisCache(cfg.RedisAddr(), cfg.RedisPassword(), logger), logger)
	service := servicepkg.NewService(repository, accessTokenConverter, refreshTokenConverter)
	handler := handlerpkg.NewHandler(cfg.IsProduction(), service, accessTokenConverter, refreshTokenConverter, logger)
	server := serverpkg.NewHTTPServer(handler, logger)

	server.MustStart(fmt.Sprintf("%s:%d", cfg.Host(), cfg.Port()), cfg.IsProduction())
}
