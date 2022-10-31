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
    step_number int NOT NULL,

    CONSTRAINT transaction_steps_transactions_fk
      FOREIGN KEY (transaction_id, transaction_version) REFERENCES transactions(id, "version")
);

COMMIT;
