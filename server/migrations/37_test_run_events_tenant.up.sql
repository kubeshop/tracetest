BEGIN;

ALTER TABLE
  test_run_events
ADD
  COLUMN tenant_id varchar
;

CREATE INDEX idx_test_run_events_tenant_id ON test_run_events(tenant_id);

COMMIT;
