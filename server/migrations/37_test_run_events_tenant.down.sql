BEGIN;

ALTER TABLE
  test_run_events
DROP
  COLUMN tenant_id;

COMMIT;
