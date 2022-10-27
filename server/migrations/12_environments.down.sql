BEGIN;

ALTER TABLE
  "test_runs" DROP COLUMN environment;

DROP TABLE "environments";

COMMIT;