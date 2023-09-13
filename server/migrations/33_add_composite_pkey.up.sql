ALTER TABLE data_stores
DROP CONSTRAINT data_stores_pkey,
ADD PRIMARY KEY (id, tenant_id);