BEGIN;

ALTER TABLE
  config DROP COLUMN tenant_id;

ALTER TABLE
  data_stores DROP COLUMN tenant_id;

ALTER TABLE
  demos DROP COLUMN tenant_id;

ALTER TABLE
  environments DROP COLUMN tenant_id;

ALTER TABLE
  linters DROP COLUMN tenant_id;

ALTER TABLE
  polling_profiles DROP COLUMN tenant_id;

ALTER TABLE
  test_runners DROP COLUMN tenant_id;

ALTER TABLE
  test_runs DROP COLUMN tenant_id;

ALTER TABLE
  tests DROP COLUMN tenant_id;

ALTER TABLE
  transaction_runs DROP COLUMN tenant_id;

ALTER TABLE
  transactions DROP COLUMN tenant_id;

ALTER TABLE
  transaction_run_steps DROP COLUMN tenant_id;

ALTER TABLE
  transaction_steps DROP COLUMN tenant_id;

COMMIT;