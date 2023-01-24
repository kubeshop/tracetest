BEGIN;

  ALTER TABLE transaction_runs
    DROP COLUMN "pass",
    DROP COLUMN "fail";

COMMIT;
