ALTER TABLE tests
ADD COLUMN "version" int NOT NULL DEFAULT 1;

CREATE UNIQUE INDEX idx_unique_test_id_version ON tests (id, "version");