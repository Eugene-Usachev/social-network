package cassandra

import (
	"fmt"

	loggerpkg "github.com/Eugune-Usachev/social-network/src/pkg/logger"
	"github.com/gocql/gocql"
)

type logger struct {
	logger loggerpkg.Logger
}

var _ gocql.StdLogger = (*logger)(nil)

func newCassandraLogger(providedLogger loggerpkg.Logger) *logger {
	return &logger{
		providedLogger,
	}
}

func (c *logger) Print(v ...interface{}) {
	c.logger.Info(fmt.Sprint(v...))
}

func (c *logger) Printf(format string, v ...interface{}) {
	c.logger.Info(fmt.Sprintf(format, v...))
}

func (c *logger) Println(v ...interface{}) {
	c.logger.Info(fmt.Sprintln(v...) + "\n")
}
