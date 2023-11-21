BEGIN;

--  Tests
ALTER TABLE
  tests
ADD
  COLUMN skip_trace_collection BOOLEAN DEFAULT FALSE NOT NULL;

--  Test Runs
ALTER TABLE
  test_runs
ADD
  COLUMN skip_trace_collection BOOLEAN DEFAULT FALSE NOT NULL;

COMMIT;