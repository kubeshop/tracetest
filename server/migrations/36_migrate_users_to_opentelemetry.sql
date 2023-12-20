-- We are going to deprecate "no tracing mode" and "agent" tracing backends
-- So we need to migrate all users to a valid trace backend.
--
-- This means that all users using "no tracing mode" and "agent" tracing backends will become
-- "otlp" tracing backends

-- If there's an "agent" tracing backend, replace it with an "otlp" tracing backend instead
UPDATE data_stores
    SET "name" = 'otlp', "type" = 'otlp', "values" = '{}'::jsonb
from (
    SELECT id, "type" FROM data_stores WHERE "type" = 'agent'
) migration_target
WHERE data_stores.id = migration_target.id;

-- If there's no "current" tracing backend, add one for otlp. This ensures that if user is on
-- "No tracing mode", it will be migrated to "OpenTelemetry".
INSERT INTO
    data_stores (id, "name", "type", is_default, "values", created_at, tenant_id)
VALUES ('current', 'otlp', 'otlp', true, '{}'::jsonb, now(), '')
ON CONFLICT DO NOTHING;


