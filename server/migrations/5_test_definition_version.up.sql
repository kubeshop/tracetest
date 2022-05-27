ALTER TABLE definitions
ADD COLUMN test_version int not null default 1;

ALTER TABLE definitions
ADD CONSTRAINT definition_test_id_version_fk
FOREIGN KEY (test_id, test_version) REFERENCES tests(id, "version");

ALTER TABLE definitions
DROP CONSTRAINT definitions_pkey;

ALTER TABLE definitions
ADD CONSTRAINT definitions_pkey PRIMARY KEY (test_id, test_version);