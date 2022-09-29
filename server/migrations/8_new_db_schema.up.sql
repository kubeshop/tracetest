BEGIN;

  ALTER TABLE runs DROP CONSTRAINT test_id_version_fk;
  DROP TABLE runs;
  DROP TABLE tests;


  CREATE TABLE "tests" (
    "id" varchar,
    "version" int,
    "name" varchar,
    "description" varchar,
    "service_under_test" jsonb,
    "specs" jsonb,
    "created_at" timestamp,
    PRIMARY KEY ("id", "version")
  );

  CREATE TABLE "test_runs" (
    "id" SERIAL PRIMARY KEY,
    "test_id" varchar,
    "test_version" int,

    -- timestamps
    "created_at" timestamp,
    "service_triggered_at" timestamp,
    "service_trigger_completed_at" timestamp,
    "obtained_trace_at" timestamp,
    "completed_at" timestamp,

    -- trigger params
    "state" varchar,
    "trace_id" varchar,
    "span_id" varchar,

    -- result info
    "trigger_results" jsonb,
    "test_results" jsonb,
    "trace" jsonb,
    "last_error" varchar,

    "metadata" jsonb,

    CONSTRAINT fk_test_runs_tests
      FOREIGN KEY ("test_id", "test_version") REFERENCES "tests" ("id", "version")
  );

COMMIT;
