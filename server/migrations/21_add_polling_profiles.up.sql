BEGIN;

CREATE TABLE "polling_profiles" (
  "id" varchar not null primary key,
  "name" varchar not null,
  "default" boolean not null default false,
  "strategy" varchar not null,
  "periodic" JSONB
);

COMMIT;
