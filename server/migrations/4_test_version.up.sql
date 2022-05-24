ALTER TABLE tests
ADD COLUMN "version" int NOT NULL DEFAULT 1;

CREATE UNIQUE INDEX idx_unique_test_id_version ON tests (id, "version");

ALTER TABLE runs
ADD COLUMN test_version int NOT NULL;

ALTER TABLE runs
ADD CONSTRAINT test_id_version_fk
FOREIGN KEY (test_id, test_version) REFERENCES tests(id, "version");