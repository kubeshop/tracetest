BEGIN;

ALTER TABLE transaction_runs
  DROP COLUMN steps,
  DROP COLUMN step_runs
;

CREATE TABLE transaction_run_steps (
    transaction_run_id VARCHAR NOT NULL,
    transaction_run_transaction_id VARCHAR NOT NULL,
    test_run_id int NOT NULL,
    test_run_test_id VARCHAR NOT NULL,

    CONSTRAINT transaction_run_steps_transaction_runs_fk
      FOREIGN KEY (transaction_run_id, transaction_run_transaction_id) REFERENCES transaction_runs(id, transaction_id),
    CONSTRAINT transaction_run_steps_test_runs_fk
      FOREIGN KEY (test_run_id, test_run_test_id) REFERENCES test_runs(id, test_id)
);


COMMIT;
