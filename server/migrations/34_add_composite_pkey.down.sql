BEGIN;

-- Tests
ALTER TABLE
  tests DROP CONSTRAINT tests_pkey CASCADE,
ADD
  CONSTRAINT tests_pkey PRIMARY KEY (id, version),
ALTER COLUMN
  tenant_id DROP DEFAULT,
ALTER COLUMN
  tenant_id DROP NOT NULL;

UPDATE
  tests
SET
  tenant_id = null
WHERE
  tenant_id = '';

ALTER TABLE
  tests
ALTER COLUMN
  tenant_id TYPE uuid using tenant_id :: uuid;

-- # Test Runs
ALTER TABLE
  test_runs DROP CONSTRAINT test_runs_pkey CASCADE,
ADD
  CONSTRAINT test_runs_pkey PRIMARY KEY (id, test_id),
ALTER COLUMN
  tenant_id DROP DEFAULT,
ALTER COLUMN
  tenant_id DROP NOT NULL;

UPDATE
  test_runs
SET
  tenant_id = null
WHERE
  tenant_id = '';

ALTER TABLE
  test_runs
ALTER COLUMN
  tenant_id TYPE uuid using tenant_id :: uuid;

-- Test Suites
ALTER TABLE
  test_suites DROP CONSTRAINT transaction_pkey CASCADE,
ADD
  CONSTRAINT transaction_pkey PRIMARY KEY (id, version),
ALTER COLUMN
  tenant_id DROP DEFAULT,
ALTER COLUMN
  tenant_id DROP NOT NULL;

UPDATE
  test_suites
SET
  tenant_id = null
WHERE
  tenant_id = '';

ALTER TABLE
  test_suites
ALTER COLUMN
  tenant_id TYPE uuid using tenant_id :: uuid;

-- Test Suite Runs
ALTER TABLE
  test_suite_runs DROP CONSTRAINT transaction_run_pkey CASCADE,
ADD
  CONSTRAINT transaction_run_pkey PRIMARY KEY (id, test_suite_id),
ALTER COLUMN
  tenant_id DROP DEFAULT,
ALTER COLUMN
  tenant_id DROP NOT NULL;

UPDATE
  test_suite_runs
SET
  tenant_id = null
WHERE
  tenant_id = '';

ALTER TABLE
  test_suite_runs
ALTER COLUMN
  tenant_id TYPE uuid using tenant_id :: uuid;

-- Variable Sets
ALTER TABLE
  variable_sets DROP CONSTRAINT environments_pkey CASCADE,
ADD
  CONSTRAINT environments_pkey PRIMARY KEY (id),
ALTER COLUMN
  tenant_id DROP DEFAULT,
ALTER COLUMN
  tenant_id DROP NOT NULL;

UPDATE
  variable_sets
SET
  tenant_id = null
WHERE
  tenant_id = '';

ALTER TABLE
  variable_sets
ALTER COLUMN
  tenant_id TYPE uuid using tenant_id :: uuid;

COMMIT;