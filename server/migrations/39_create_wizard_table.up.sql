BEGIN;

CREATE TABLE "wizards" (
  "tenant_id" varchar PRIMARY KEY NOT NULL DEFAULT '',
  "steps" JSONB
);

COMMIT;
