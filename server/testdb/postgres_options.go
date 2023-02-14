package testdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/j2gg0s/otsql"
	"github.com/j2gg0s/otsql/hook/trace"
	pq "github.com/lib/pq"
	"go.opentelemetry.io/otel/attribute"
)

type PostgresOption func(*postgresDB) error

func dbSpanNameFormatter(ctx context.Context, method, query string) string {
	splitQuery := strings.Fields(query)
	queryName := ""
	if len(splitQuery) > 0 {
		queryName = splitQuery[0]
	}

	queryName = strings.ReplaceAll(queryName, "\n", "")

	return fmt.Sprintf("%s %s", method, queryName)
}

func Connect(dsn string) (*sql.DB, error) {
	connector, err := pq.NewConnector(dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}
	db := sql.OpenDB(
		otsql.WrapConnector(connector,
			otsql.WithHooks(
				trace.New(
					trace.WithQuery(true),
					trace.WithQueryParams(true),
					trace.WithRowsAffected(true),
					trace.WithSpanNameFormatter(dbSpanNameFormatter),
					trace.WithDefaultAttributes(attribute.String("service.name", "tracetest")),
				),
			),
		),
	)

	return db, nil
}

func WithDB(db *sql.DB) PostgresOption {
	return func(pd *postgresDB) error {
		pd.db = db
		return nil
	}
}
