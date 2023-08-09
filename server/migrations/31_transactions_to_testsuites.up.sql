BEGIN;

ALTER TABLE
  transactions RENAME TO test_suites;

ALTER TABLE
  transaction_steps RENAME TO test_suite_steps;

ALTER TABLE
  test_suite_steps RENAME COlUMN transaction_id TO test_suite_id;

ALTER TABLE
  test_suite_steps RENAME COlUMN transaction_version TO test_suite_version;

ALTER TABLE
  transaction_runs RENAME TO test_suite_runs;

ALTER TABLE
  test_suite_runs RENAME COlUMN transaction_id TO test_suite_id;

ALTER TABLE
  test_suite_runs RENAME COlUMN transaction_version TO test_suite_version;

ALTER TABLE
  transaction_run_steps RENAME TO test_suite_run_steps;

ALTER TABLE
  test_suite_run_steps RENAME COLUMN transaction_run_id TO test_suite_run_id;

ALTER TABLE
  test_suite_run_steps RENAME COLUMN transaction_run_transaction_id TO test_suite_run_test_suite_id;

COMMIT;