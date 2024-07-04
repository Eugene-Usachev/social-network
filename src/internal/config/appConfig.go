package config

import (
	fb "github.com/Eugene-Usachev/fastbytes"
	"github.com/goccy/go-json"
	"log"
	"os"
	"strconv"
)

type AppConfig struct {
	isProduction bool
	port         int
	host         string

	esAddr []string
	esUser string
	esPass string

	postgresHost    string
	postgresPort    int
	postgresUser    string
	postgresPass    string
	postgresDBName  string
	postgresSSLMode string

	redisAddr     string
	redisPassword string

	metricsAddr string

	fstAccessKey  string
	fstRefreshKey string
}

func MustNewConfig() *AppConfig {
	c := &AppConfig{}

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

	host := os.Getenv("HOST")
	if host != "" {
		c.host = host
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

	postgresSSLMode := os.Getenv("POSTGRES_SSL_MODE")
	if postgresSSLMode != "" {
		c.postgresSSLMode = postgresSSLMode
	} else {
		log.Fatal("POSTGRES_SSL_MODE is not set")
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

	metricsAddr := os.Getenv("METRICS_ADDRESS")
	if metricsAddr != "" {
		c.metricsAddr = metricsAddr
	} else {
		log.Fatal("METRICS_ADDRESS is not set")
	}

	return c
}

func (appConfig *AppConfig) IsProduction() bool {
	return appConfig.isProduction
}

func (appConfig *AppConfig) Port() int {
	return appConfig.port
}

func (appConfig *AppConfig) Host() string {
	return appConfig.host
}

func (appConfig *AppConfig) EsAddr() []string {
	return appConfig.esAddr
}

func (appConfig *AppConfig) EsUser() string {
	return appConfig.esUser
}

func (appConfig *AppConfig) EsPass() string {
	return appConfig.esPass
}

func (appConfig *AppConfig) PostgresHost() string {
	return appConfig.postgresHost
}

func (appConfig *AppConfig) PostgresPort() int {
	return appConfig.postgresPort
}

func (appConfig *AppConfig) PostgresUser() string {
	return appConfig.postgresUser
}

func (appConfig *AppConfig) PostgresPass() string {
	return appConfig.postgresPass
}

func (appConfig *AppConfig) PostgresDBName() string {
	return appConfig.postgresDBName
}

func (appConfig *AppConfig) PostgresSSLMode() string { return appConfig.postgresSSLMode }

func (appConfig *AppConfig) RedisAddr() string {
	return appConfig.redisAddr
}

func (appConfig *AppConfig) RedisPassword() string {
	return appConfig.redisPassword
}

func (appConfig *AppConfig) MetricsAddr() string {
	return appConfig.metricsAddr
}

func (appConfig *AppConfig) FstAccessKey() string {
	return appConfig.fstAccessKey
}

func (appConfig *AppConfig) FstRefreshKey() string {
	return appConfig.fstRefreshKey
}
