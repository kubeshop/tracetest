BEGIN;

ALTER TABLE "test_run_events" ADD COLUMN "title" varchar not null;

COMMIT;
