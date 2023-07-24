BEGIN;

ALTER TABLE
  config
ADD
  COLUMN tenant_id uuid;

CREATE INDEX idx_config_tenant_id ON config(tenant_id);

ALTER TABLE
  data_stores
ADD
  COLUMN tenant_id uuid;

CREATE INDEX idx_data_stores_tenant_id ON data_stores(tenant_id);

ALTER TABLE
  demos
ADD
  COLUMN tenant_id uuid;

CREATE INDEX idx_demos_tenant_id ON demos(tenant_id);

ALTER TABLE
  environments
ADD
  COLUMN tenant_id uuid;

CREATE INDEX idx_environments_tenant_id ON environments(tenant_id);

ALTER TABLE
  linters
ADD
  COLUMN tenant_id uuid;

CREATE INDEX idx_linters_tenant_id ON linters(tenant_id);

ALTER TABLE
  polling_profiles
ADD
  COLUMN tenant_id uuid;

CREATE INDEX idx_polling_profiles_tenant_id ON polling_profiles(tenant_id);

ALTER TABLE
  test_runners
ADD
  COLUMN tenant_id uuid;

CREATE INDEX idx_test_runners_tenant_id ON test_runners(tenant_id);

ALTER TABLE
  test_runs
ADD
  COLUMN tenant_id uuid;

CREATE INDEX idx_test_runs_tenant_id ON test_runs(tenant_id);

ALTER TABLE
  tests
ADD
  COLUMN tenant_id uuid;

CREATE INDEX idx_tests_tenant_id ON tests(tenant_id);

ALTER TABLE
  transaction_run_steps
ADD
  COLUMN tenant_id uuid;

CREATE INDEX idx_transaction_run_steps_tenant_id ON transaction_run_steps(tenant_id);

ALTER TABLE
  transaction_runs
ADD
  COLUMN tenant_id uuid;

CREATE INDEX idx_transaction_runs_tenant_id ON transaction_runs(tenant_id);

ALTER TABLE
  transaction_steps
ADD
  COLUMN tenant_id uuid;

CREATE INDEX idx_transaction_steps_tenant_id ON transaction_steps(tenant_id);

ALTER TABLE
  transactions
ADD
  COLUMN tenant_id uuid;

CREATE INDEX idx_transactions_tenant_id ON transactions(tenant_id);

COMMIT;