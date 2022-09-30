BEGIN;
  ALTER TABLE "test_runs" DROP CONSTRAINT fk_test_runs_tests;
  DROP TABLE "test_runs";
  DROP TABLE "tests";

  CREATE TABLE runs (
      id uuid NOT NULL,
      test_id uuid NOT NULL,
      run json NOT NULL,
      test_version integer DEFAULT 1 NOT NULL
  );

  CREATE TABLE tests (
      id uuid NOT NULL,
      test json NOT NULL,
      version integer DEFAULT 1 NOT NULL
  );

  ALTER TABLE ONLY runs
      ADD CONSTRAINT runs_pkey PRIMARY KEY (id);

  ALTER TABLE ONLY tests
      ADD CONSTRAINT tests_pkey PRIMARY KEY (id, version);


  ALTER TABLE ONLY runs
      ADD CONSTRAINT test_id_version_fk FOREIGN KEY (test_id, test_version) REFERENCES tests(id, version);

COMMIT;
