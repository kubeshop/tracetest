BEGIN;

--  Tests
ALTER TABLE
  tests
ALTER COLUMN
  tenant_id TYPE varchar;

UPDATE
  tests
SET
  tenant_id = ''
WHERE
  tenant_id is null;

ALTER TABLE
  tests DROP CONSTRAINT tests_pkey CASCADE,
ADD
  CONSTRAINT tests_pkey PRIMARY KEY (id, version, tenant_id),
ALTER COLUMN
  tenant_id
SET
  DEFAULT '';

--  Test Runs
ALTER TABLE
  test_runs
ALTER COLUMN
  tenant_id TYPE varchar;

UPDATE
  test_runs
SET
  tenant_id = ''
WHERE
  tenant_id is null;

ALTER TABLE
  test_runs DROP CONSTRAINT test_runs_pkey CASCADE,
ADD
  CONSTRAINT test_runs_pkey PRIMARY KEY (id, test_id, tenant_id),
ADD
  CONSTRAINT fk_test_runs_tests FOREIGN KEY (test_id, test_version, tenant_id) REFERENCES tests(id, version, tenant_id),
ALTER COLUMN
  tenant_id
SET
  DEFAULT '';

--  Test Suites
ALTER TABLE
  test_suites
ALTER COLUMN
  tenant_id TYPE varchar;

UPDATE
  test_suites
SET
  tenant_id = ''
WHERE
  tenant_id is null;

ALTER TABLE
  test_suites DROP CONSTRAINT transaction_pkey CASCADE,
ADD
  CONSTRAINT transaction_pkey PRIMARY KEY (id, version, tenant_id),
ALTER COLUMN
  tenant_id
SET
  DEFAULT '';

--  Test Suite Runs
ALTER TABLE
  test_suite_runs
ALTER COLUMN
  tenant_id TYPE varchar;

UPDATE
  test_suite_runs
SET
  tenant_id = ''
WHERE
  tenant_id is null;

ALTER TABLE
  test_suite_runs DROP CONSTRAINT transaction_run_pkey CASCADE,
ADD
  CONSTRAINT transaction_run_pkey PRIMARY KEY (id, test_suite_id, tenant_id),
ADD
  CONSTRAINT transaction_run_transactions_fk FOREIGN KEY (test_suite_id, test_suite_version, tenant_id) REFERENCES test_suites(id, version, tenant_id),
ALTER COLUMN
  tenant_id
SET
  DEFAULT '';

--  Variable Sets
ALTER TABLE
  variable_sets
ALTER COLUMN
  tenant_id TYPE varchar;

UPDATE
  variable_sets
SET
  tenant_id = ''
WHERE
  tenant_id is null;

ALTER TABLE
  variable_sets DROP CONSTRAINT environments_pkey CASCADE,
ADD
  CONSTRAINT environments_pkey PRIMARY KEY (id, tenant_id),
ALTER COLUMN
  tenant_id
SET
  DEFAULT '';

COMMIT;