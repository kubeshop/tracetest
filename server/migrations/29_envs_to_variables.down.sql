BEGIN;

ALTER TABLE
  variable_sets RENAME TO environments;

ALTER TABLE
  test_runs RENAME COLUMN variable_set TO environment;

ALTER TABLE
  transaction_runs RENAME COLUMN variable_set TO environment;

COMMIT;