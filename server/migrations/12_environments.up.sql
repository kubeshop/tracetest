BEGIN;

CREATE TABLE "environments" (
  "id" varchar not null primary key,
  "name" varchar,
  "description" varchar,
  "values" jsonb,
  "created_at" timestamp
);

ALTER TABLE
  "test_runs"
ADD
  COLUMN environment jsonb;

COMMIT;