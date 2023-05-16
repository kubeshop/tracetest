BEGIN;

CREATE TABLE "linterns" (
  "id" varchar not null primary key,
  "name" varchar not null,
  "enabled" boolean not null default true,
  "minimumScore" integer not null,
  "plugins" JSONB
);

ALTER TABLE test_runs
    ADD COLUMN lintern JSONB;

COMMIT;
