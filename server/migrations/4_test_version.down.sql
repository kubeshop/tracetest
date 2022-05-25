ALTER TABLE runs
DROP CONSTRAINT test_id_version_fk;

ALTER TABLE runs
DROP COLUMN test_version;

ALTER TABLE tests
DROP CONSTRAINT tests_pkey;

ALTER TABLE tests
ADD CONSTRAINT tests_pkey PRIMARY KEY (id);

ALTER TABLE tests
DROP COLUMN "version";