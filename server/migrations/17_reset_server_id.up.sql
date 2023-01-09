-- Before this version, all tracetest versions had a hard-coded machine-id
-- this means that the id stored in this table changed by version, but two
-- users running the same version would share the same id.
--
-- With the problem fixed, we need to make sure we reset the server id so all
-- users have real unique ids.
DELETE FROM "server";
