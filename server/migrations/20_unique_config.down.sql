BEGIN;

DROP TABLE "config";

CREATE TABLE "configs" (
  "id" varchar not null primary key,
  "name" varchar,
  "analytics_enabled" boolean not null
);

COMMIT;
