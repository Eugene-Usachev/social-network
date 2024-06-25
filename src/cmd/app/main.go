package main

import (
	"fmt"
	fb "github.com/Eugene-Usachev/fastbytes"
	"github.com/goccy/go-json"
	"log"
	"os"
	loggerpkg "social-network/src/pkg/logger"
	"strconv"
	"time"
)

type config struct {
	isProduction bool
	port         int

	esAddr []string
	esUser string
	esPass string

	postgresHost   string
	postgresPort   int
	postgresUser   string
	postgresPass   string
	postgresDBName string

	redisAddr     string
	redisPassword string

	prometheusAddr string

	fstAccessKey  string
	fstRefreshKey string
}

func getConfig() *config {
	c := &config{}

	isProduction := os.Getenv("IS_PRODUCTION")
	if isProduction != "" {
		c.isProduction, _ = strconv.ParseBool(isProduction)
	} else {
		log.Fatal("IS_PRODUCTION is not set")
	}

	port := os.Getenv("PORT")
	if port != "" {
		var err error
		c.port, err = strconv.Atoi(port)
		if err != nil {
			log.Fatal("Failed to parse PORT: ", err)
		}
	} else {
		log.Fatal("PORT is not set")
	}

	esAddresses := os.Getenv("ES_ADDRESSES")
	if esAddresses != "" {
		var addresses []string
		if err := json.Unmarshal(fb.S2B(esAddresses), &addresses); err != nil {
			log.Fatal("Failed to unmarshal ES_ADDRESSES: ", err)
		}
		c.esAddr = addresses
	} else {
		log.Fatal("ES_ADDR is not set")
	}

	esUser := os.Getenv("ES_USERNAME")
	if esUser != "" {
		c.esUser = esUser
	} else {
		log.Fatal("ES_USERNAME is not set")
	}

	esPass := os.Getenv("ES_PASSWORD")
	if esPass != "" {
		c.esPass = esPass
	} else {
		log.Fatal("ES_PASSWORD is not set")
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost != "" {
		c.postgresHost = postgresHost
	} else {
		log.Fatal("POSTGRES_HOST is not set")
	}

	postgresPort := os.Getenv("POSTGRES_PORT")
	if postgresPort != "" {
		var err error
		c.postgresPort, err = strconv.Atoi(postgresPort)
		if err != nil {
			log.Fatal("POSTGRES_PORT is not a number")
		}
	} else {
		log.Fatal("POSTGRES_PORT is not set")
	}

	postgresUser := os.Getenv("POSTGRES_USERNAME")
	if postgresUser != "" {
		c.postgresUser = postgresUser
	} else {
		log.Fatal("POSTGRES_USERNAME is not set")
	}

	postgresPass := os.Getenv("POSTGRES_PASSWORD")
	if postgresPass != "" {
		c.postgresPass = postgresPass
	} else {
		log.Fatal("POSTGRES_PASSWORD is not set")
	}

	postgresDBName := os.Getenv("POSTGRES_DB_NAME")
	if postgresDBName != "" {
		c.postgresDBName = postgresDBName
	} else {
		log.Fatal("POSTGRES_DATABASE is not set")
	}

	redisAddr := os.Getenv("REDIS_ADDRESS")
	if redisAddr != "" {
		c.redisAddr = redisAddr
	} else {
		log.Fatal("REDIS_ADDRESS is not set")
	}

	redisPort := os.Getenv("REDIS_PASSWORD")
	if redisPort != "" {
		c.redisPassword = redisPort
	} else {
		log.Fatal("REDIS_PORT is not set")
	}

	fstAccessKey := os.Getenv("FST_ACCESS_KEY")
	if fstAccessKey != "" {
		c.fstAccessKey = fstAccessKey
	} else {
		log.Fatal("FST_ACCESS_KEY is not set")
	}

	fstRefreshKey := os.Getenv("FST_REFRESH_KEY")
	if fstRefreshKey != "" {
		c.fstRefreshKey = fstRefreshKey
	} else {
		log.Fatal("FST_REFRESH_KEY is not set")
	}

	prometheusAddr := os.Getenv("PROMETHEUS_ADDRESS")
	if prometheusAddr != "" {
		c.prometheusAddr = prometheusAddr
	} else {
		log.Fatal("PROMETHEUS_ADDRESS is not set")
	}

	return c
}

func main() {
	cfg := getConfig()

	var logger loggerpkg.Logger

	if cfg.isProduction {
		logger = loggerpkg.NewElasticSearchLogger(cfg.esAddr, cfg.esUser, cfg.esPass)
	} else {
		logger = loggerpkg.NewZeroLogger()
	}

	for i := 0; i < 1000; i++ {
		logger.Info(fmt.Sprintf("test: %d", i))
	}

	time.Sleep(time.Hour)
}
