CREATE TABLE definitions (
    test_id uuid NOT NULL,
    definition json NOT NULL,
    test_version integer DEFAULT 1 NOT NULL
);

ALTER TABLE definitions
ADD CONSTRAINT definition_test_id_version_fk
FOREIGN KEY (test_id, test_version) REFERENCES tests(id, "version");

ALTER TABLE definitions
ADD CONSTRAINT definitions_pkey PRIMARY KEY (test_id, test_version);
