BEGIN;

ALTER TABLE
  environments RENAME TO variable_sets;

ALTER TABLE
  test_runs RENAME COLUMN environment TO variable_set;

ALTER TABLE
  transaction_runs RENAME COLUMN environment TO variable_set;

COMMIT;