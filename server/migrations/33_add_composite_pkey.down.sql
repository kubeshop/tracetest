BEGIN;

ALTER TABLE data_stores
DROP CONSTRAINT data_stores_pkey,
ADD PRIMARY KEY (id),
ALTER COLUMN tenant_id DROP DEFAULT,
ALTER COLUMN tenant_id DROP NOT NULL;

UPDATE data_stores
SET tenant_id = null
WHERE tenant_id = '';

ALTER TABLE data_stores ALTER COLUMN tenant_id TYPE uuid using tenant_id::uuid;


ALTER TABLE demos
DROP CONSTRAINT demos_pkey,
ADD PRIMARY KEY (id),
ALTER COLUMN tenant_id DROP DEFAULT,
ALTER COLUMN tenant_id DROP NOT NULL;

UPDATE demos
SET tenant_id = null
WHERE tenant_id = '';

ALTER TABLE demos ALTER COLUMN tenant_id TYPE uuid using tenant_id::uuid;


ALTER TABLE polling_profiles
DROP CONSTRAINT polling_profiles_pkey,
ADD PRIMARY KEY (id),
ALTER COLUMN tenant_id DROP DEFAULT,
ALTER COLUMN tenant_id DROP NOT NULL;

UPDATE polling_profiles
SET tenant_id = null
WHERE tenant_id = '';

ALTER TABLE polling_profiles ALTER COLUMN tenant_id TYPE uuid using tenant_id::uuid;


ALTER TABLE linters
DROP CONSTRAINT linters_pkey,
ADD PRIMARY KEY (id),
ALTER COLUMN tenant_id DROP DEFAULT,
ALTER COLUMN tenant_id DROP NOT NULL;

UPDATE linters
SET tenant_id = null
WHERE tenant_id = '';

ALTER TABLE linters ALTER COLUMN tenant_id TYPE uuid using tenant_id::uuid;


ALTER TABLE test_runners
DROP CONSTRAINT test_runners_pkey,
ADD PRIMARY KEY (id),
ALTER COLUMN tenant_id DROP DEFAULT,
ALTER COLUMN tenant_id DROP NOT NULL;

UPDATE test_runners
SET tenant_id = null
WHERE tenant_id = '';

ALTER TABLE test_runners ALTER COLUMN tenant_id TYPE uuid using tenant_id::uuid;

COMMIT;