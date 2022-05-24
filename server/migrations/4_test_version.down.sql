ALTER TABLE tests
DROP COLUMN "version";

DROP INDEX idx_unique_test_id_version;


ALTER TABLE runs
DROP COLUMN test_version;

ALTER TABLE runs
DROP CONSTRAINT test_id_version_fk;