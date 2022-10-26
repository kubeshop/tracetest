BEGIN;

ALTER TABLE
  "test_runs" DROP COLUMN env_snapshot;

DROP TABLE "environments";

COMMIT;