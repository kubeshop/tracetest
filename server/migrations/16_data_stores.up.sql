BEGIN;

CREATE TABLE "data_stores" (
  "id" varchar not null primary key,
  "name" varchar,
  "type" varchar not null,
  "is_default" boolean not null,
  "values" jsonb,
  "created_at" timestamp
);

COMMIT;
