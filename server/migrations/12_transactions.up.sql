BEGIN;

CREATE TABLE transactions (
    id VARCHAR NOT NULL,
    version integer NOT NULL,
    name varchar,
    description varchar,
    created_at timestamp,

    CONSTRAINT transaction_pkey PRIMARY KEY (id, version)
);

CREATE TABLE transaction_steps (
    transaction_id varchar NOT NULL,
    transaction_version integer NOT NULL,
    test_id varchar NOT NULL,
    test_version integer  NOT NULL,

    CONSTRAINT transaction_steps_transactions_fk
      FOREIGN KEY (transaction_id, transaction_version) REFERENCES transactions(id, "version"),
    CONSTRAINT transaction_steps_tests_fk
      FOREIGN KEY (test_id, test_version) REFERENCES tests(id, "version")
);

COMMIT;
