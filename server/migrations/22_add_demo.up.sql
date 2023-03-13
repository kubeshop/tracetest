BEGIN;

CREATE TABLE "demos" (
  "id" varchar not null primary key,
  "name" varchar not null,
  "type" varchar not null,
  "enabled" boolean not null,
  "pokeshop" JSONB,
  "opentelemetry_store" JSONB
);

COMMIT;
