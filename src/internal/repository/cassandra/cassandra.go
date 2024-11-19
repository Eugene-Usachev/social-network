package cassandra

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	loggerpkg "github.com/Eugune-Usachev/social-network/src/pkg/logger"
	"github.com/Eugune-Usachev/social-network/src/pkg/utils"
	"github.com/gocql/gocql"
)

type Config struct {
	Hosts    []string
	Port     int
	Keyspace string
}

func MustCreateCassandraClient(cfg Config, logger loggerpkg.Logger) *gocql.Session {
	var session *gocql.Session
	var err error

	err = utils.DoWithTries(func() error {
		cluster := gocql.NewCluster(cfg.Hosts...)
		cluster.Port = cfg.Port
		cluster.Keyspace = cfg.Keyspace
		cluster.Logger = newCassandraLogger(logger)

		session, err = cluster.CreateSession()
		if err != nil {
			return err
		}

		// Check connection by querying system tables
		if err = session.Query("SELECT release_version FROM system.local").Exec(); err != nil {
			return err
		}

		return nil
	}, 20, 1*time.Second)

	if err != nil {
		logger.Fatal(fmt.Sprintf("Error connecting to Cassandra: %s", err.Error()))
	}

	logger.Info("Creating tables for Cassandra")

	err = createTablesAndIndexes(session)
	if err != nil {
		logger.Fatal(err.Error())
	} else {
		logger.Info("All tables for Cassandra created")
	}

	return session
}

//go:embed queries/up.cql
var upQuery string

func createTablesAndIndexes(session *gocql.Session) error {
	// Split queries by semicolon for execution
	queries := splitCQLQueries(upQuery)
	for _, query := range queries {
		if query == "" {
			continue
		}

		if err := session.Query(query).Exec(); err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %w", query, err)
		}
	}

	return nil
}

func splitCQLQueries(cql string) []string {
	// Split by semicolons for multiple queries in the CQL file
	// Adjust logic if the queries span multiple lines
	queries := strings.Split(cql, ";")
	for i := range queries {
		queries[i] = strings.TrimSpace(queries[i])
	}
	return queries
}
