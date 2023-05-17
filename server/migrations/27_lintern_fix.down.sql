BEGIN;

ALTER TABLE "linterns" RENAME COLUMN "minimum_score" to "minimumScore";

COMMIT;
