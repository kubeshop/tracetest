BEGIN;

ALTER TABLE data_stores ALTER COLUMN tenant_id TYPE varchar;

UPDATE data_stores
SET tenant_id = ''
WHERE tenant_id is null;

ALTER TABLE data_stores
DROP CONSTRAINT data_stores_pkey,
ADD PRIMARY KEY (id, tenant_id),
ALTER COLUMN tenant_id SET DEFAULT '';


ALTER TABLE demos ALTER COLUMN tenant_id TYPE varchar;

UPDATE demos
SET tenant_id = ''
WHERE tenant_id is null;

ALTER TABLE demos
DROP CONSTRAINT demos_pkey,
ADD PRIMARY KEY (id, tenant_id),
ALTER COLUMN tenant_id SET DEFAULT '';


ALTER TABLE polling_profiles ALTER COLUMN tenant_id TYPE varchar;

UPDATE polling_profiles
SET tenant_id = ''
WHERE tenant_id is null;

ALTER TABLE polling_profiles
DROP CONSTRAINT polling_profiles_pkey,
ADD PRIMARY KEY (id, tenant_id),
ALTER COLUMN tenant_id SET DEFAULT '';


ALTER TABLE linters ALTER COLUMN tenant_id TYPE varchar;

UPDATE linters
SET tenant_id = ''
WHERE tenant_id is null;

ALTER TABLE linters
DROP CONSTRAINT linters_pkey,
ADD PRIMARY KEY (id, tenant_id),
ALTER COLUMN tenant_id SET DEFAULT '';


ALTER TABLE test_runners ALTER COLUMN tenant_id TYPE varchar;

UPDATE test_runners
SET tenant_id = ''
WHERE tenant_id is null;

ALTER TABLE test_runners
DROP CONSTRAINT test_runners_pkey,
ADD PRIMARY KEY (id, tenant_id),
ALTER COLUMN tenant_id SET DEFAULT '';

COMMIT;