BEGIN;

ALTER TABLE transaction_runs
  ADD COLUMN steps jsonb NOT NULL,
  ADD COLUMN  step_runs jsonb NOT NULL
;

DROP TABLE transaction_run_steps;

COMMIT;
