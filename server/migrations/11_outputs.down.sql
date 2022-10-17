BEGIN;

  ALTER TABLE tests DROP COLUMN outputs;
  ALTER TABLE test_runs DROP COLUMN outputs;

COMMIT;
