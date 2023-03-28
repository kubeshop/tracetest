BEGIN;

CREATE TABLE "test_run_events" (
  "id" SERIAL PRIMARY KEY,
  "test_id" varchar not null,
  "run_id" int not null,
  "type" varchar not null,
  "stage" varchar not null,
  "description" varchar not null,
  "created_at" timestamp not null default now(),
  "data_store_connection" JSONB,
  "polling" JSONB,
  "outputs" JSONB
);

CREATE INDEX test_run_event_test_id_idx ON test_run_events(test_id);
CREATE INDEX test_run_event_run_id_idx ON test_run_events(run_id);

COMMIT;
