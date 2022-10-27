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
  COLUMN env_snapshot jsonb;

ALTER TABLE
  "test_runs"
ADD
  COLUMN env_id varchar;

COMMIT;