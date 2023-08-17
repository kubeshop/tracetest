package traces

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/kubeshop/tracetest/server/http/middleware"
	"github.com/kubeshop/tracetest/server/pkg/sqlutil"
	"go.opentelemetry.io/otel/trace"
)

type TraceRepository struct {
	db  *sql.DB
	log func(string, ...interface{})
}

func NewTraceRepository(db *sql.DB) *TraceRepository {
	return &TraceRepository{
		log: func(format string, args ...interface{}) {
			log.Printf("[TraceRepository] "+format, args...)
		},
		db: db,
	}
}

func (r *TraceRepository) Get(ctx context.Context, id trace.TraceID) (Trace, error) {
	sql, params := sqlutil.Tenant(
		ctx,
		`SELECT trace FROM otlp_traces WHERE trace_id = $1`,
		id.String(),
	)

	var jsonTrace []byte
	err := r.db.QueryRowContext(ctx, sql, params...).Scan(&jsonTrace)
	if err != nil {
		return Trace{}, err
	}

	var t Trace
	err = json.Unmarshal(jsonTrace, &t)
	if err != nil {
		return Trace{}, fmt.Errorf("failed to unmarshal trace: %w", err)
	}

	return t, nil
}

func (r *TraceRepository) UpdateTraceSpans(ctx context.Context, trace *Trace) error {
	r.log("updating trace %s with %d spans", trace.ID.String(), len(trace.Spans()))
	old, err := r.Get(ctx, trace.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to get trace: %w", err)
	}
	r.log("got old trace %s with %d spans", old.ID.String(), len(old.Spans()))

	// we need to merge preexisting spans with the new spans
	updatedTrace := NewTrace(
		trace.ID.String(),
		append(old.Spans(), trace.Spans()...),
	)
	r.log("updated trace %s with %d spans", updatedTrace.ID.String(), len(updatedTrace.Spans()))

	jsonTrace, err := updatedTrace.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal trace: %w", err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	sql, params := sqlutil.Tenant(
		ctx,
		`DELETE FROM otlp_traces WHERE trace_id = $1`,
		trace.ID.String(),
	)

	_, err = tx.ExecContext(
		ctx,
		sql,
		params...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete trace: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO otlp_traces (tenant_id, trace_id, trace) VALUES ($1, $2, $3)`,
		middleware.TenantIDFromContext(ctx),
		trace.ID.String(),
		jsonTrace,
	)

	if err != nil {
		return fmt.Errorf("failed to insert trace: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
