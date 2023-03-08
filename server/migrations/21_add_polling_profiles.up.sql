BEGIN;

CREATE TABLE "polling_profiles" (
  "id" varchar not null primary key,
  "name" varchar not null,
  "strategy" varchar not null,
  "periodic" JSONB
);

COMMIT;
