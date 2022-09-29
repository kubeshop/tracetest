BEGIN;

  DELETE FROM test_runs;

  ALTER TABLE test_runs
    DROP CONSTRAINT test_runs_pkey,
    ALTER COLUMN id DROP DEFAULT,
    ADD PRIMARY KEY ("id", "test_id");

  DROP SEQUENCE "test_runs_id_seq";

COMMIT;
