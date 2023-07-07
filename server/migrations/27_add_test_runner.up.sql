BEGIN;

CREATE TABLE "test_runners" (
  "id" varchar not null primary key,
  "name" varchar not null,
  "required_gates" JSONB
);

ALTER TABLE test_runs
    ADD COLUMN required_gates_result JSONB;

COMMIT;
