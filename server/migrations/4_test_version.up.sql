ALTER TABLE tests
ADD COLUMN "version" int NOT NULL DEFAULT 1;

ALTER TABLE tests
DROP CONSTRAINT tests_pkey;

ALTER TABLE tests
ADD CONSTRAINT tests_pkey PRIMARY KEY (id, "version");

ALTER TABLE runs
ADD COLUMN test_version int NOT NULL DEFAULT 1;

ALTER TABLE runs
ADD CONSTRAINT test_id_version_fk
FOREIGN KEY (test_id, test_version) REFERENCES tests(id, "version");