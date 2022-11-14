BEGIN;

CREATE TABLE transaction_runs (
    id VARCHAR NOT NULL,
    transaction_id VARCHAR NOT NULL,
    transaction_version integer NOT NULL,
    created_at timestamp NOT NULL,
    completed_at timestamp NOT NULL,
    state varchar NOT NULL,
    steps jsonb NOT NULL,
    step_runs jsonb NOT NULL,
    current_test int NOT NULL,
    last_error varchar,
    metadata jsonb,
    environment jsonb,

    CONSTRAINT transaction_run_pkey PRIMARY KEY (id, transaction_id),
    CONSTRAINT transaction_run_transactions_fk
      FOREIGN KEY (transaction_id, transaction_version) REFERENCES transactions(id, "version")
);

COMMIT;
