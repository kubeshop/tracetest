BEGIN;

  ALTER TABLE transaction_runs
    ADD COLUMN "pass" int,
    ADD COLUMN "fail" int;

COMMIT;
