BEGIN;

  ALTER TABLE test_runs
    ADD COLUMN "pass" int,
    ADD COLUMN "fail" int;

COMMIT;
