ALTER TABLE definitions
DROP CONSTRAINT definitions_pkey;

ALTER TABLE definitions
ADD CONSTRAINT definitions_pkey PRIMARY KEY (test_id);

ALTER TABLE definitions
DROP CONSTRAINT definition_test_id_version_fk;

ALTER TABLE definitions
DROP COLUMN test_version;
