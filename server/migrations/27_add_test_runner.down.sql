BEGIN;

DROP TABLE "test_runners";

ALTER TABLE
  test_runs DROP COLUMN required_gates_result;

ALTER TABLE
  transaction_runs DROP COLUMN all_steps_required_gates_passed;

COMMIT;