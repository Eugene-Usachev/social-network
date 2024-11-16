package postgres

import (
	"context"
	"fmt"

	"github.com/Eugune-Usachev/social-network/src/pkg/logger"
	"github.com/jackc/pgx/v5"
)

type PGLogger struct {
	logger logger.Logger
}

func NewPostgresLogger(logger logger.Logger) *PGLogger {
	return &PGLogger{
		logger,
	}
}

func (l *PGLogger) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if len(data.Args) > 0 {
		l.logger.Info(fmt.Sprintf("%s, with args: %v", data.SQL, data.Args))
	} else {
		l.logger.Info(data.SQL)
	}

	return ctx
}

func (l *PGLogger) TraceQueryEnd(_ context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		l.logger.Error("tag: " + data.CommandTag.String() + " err: " + data.Err.Error())
	}
}
