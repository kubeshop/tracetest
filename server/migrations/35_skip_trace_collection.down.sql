BEGIN;

-- Tests
ALTER TABLE
  tests DROP COLUMN skip_trace_collection;

-- Test Runs
ALTER TABLE
  test_runs DROP COLUMN skip_trace_collection;

COMMIT;