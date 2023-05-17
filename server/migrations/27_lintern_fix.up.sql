BEGIN;

ALTER TABLE "linterns" RENAME COLUMN "minimumScore" to "minimum_score";

COMMIT;
