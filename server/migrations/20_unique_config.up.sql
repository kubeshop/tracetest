BEGIN;

DROP TABLE "configs";

CREATE TABLE "config" (
  "analytics_enabled" boolean not null
);

COMMIT;
