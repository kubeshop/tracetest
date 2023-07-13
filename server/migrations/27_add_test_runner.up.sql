BEGIN;

CREATE TABLE "test_runners" (
  "id" varchar not null primary key,
  "name" varchar not null,
  "required_gates" JSONB
);

ALTER TABLE
  test_runs
ADD
  COLUMN required_gates_result JSONB;

ALTER TABLE
  transaction_runs
ADD
  COLUMN all_steps_required_gates_passed boolean;

COMMIT;