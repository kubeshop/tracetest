BEGIN;

  ALTER TABLE test_runs
    DROP COLUMN "pass",
    DROP COLUMN "fail";

COMMIT;
