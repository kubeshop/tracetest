BEGIN;

CREATE TABLE "otlp_traces" (
  "trace_id" varchar not null primary key,
  "tenant_id" uuid,
  "trace" JSONB,
  "created_at" timestamp NOT NULL DEFAULT NOW()
);

COMMIT;
