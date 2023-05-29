BEGIN;

CREATE TABLE "linters" (
  "id" varchar not null primary key,
  "name" varchar not null,
  "enabled" boolean not null default true,
  "minimum_score" integer not null,
  "plugins" JSONB
);

ALTER TABLE test_runs
    ADD COLUMN linter JSONB;

COMMIT;
