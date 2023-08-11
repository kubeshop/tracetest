BEGIN;

ALTER TABLE
  test_suites RENAME TO transactions;

ALTER TABLE
  test_suite_steps RENAME TO transaction_steps;

ALTER TABLE
  transaction_steps RENAME COlUMN test_suite_id TO transaction_id;

ALTER TABLE
  transaction_steps RENAME COlUMN test_suite_version TO transaction_version;

ALTER TABLE
  test_suite_runs RENAME TO transaction_runs;

ALTER TABLE
  transaction_runs RENAME COlUMN test_suite_id TO transaction_id;

ALTER TABLE
  transaction_runs RENAME COlUMN test_suite_version TO transaction_version;

ALTER TABLE
  test_suite_run_steps RENAME TO transaction_run_steps;

ALTER TABLE
  transaction_run_steps RENAME COLUMN test_suite_run_id TO transaction_run_id;

ALTER TABLE
  transaction_run_steps RENAME COLUMN test_suite_run_test_suite_id TO transaction_run_transaction_id;

COMMIT;