BEGIN;

UPDATE data_stores
SET id = 'current'
WHERE id = (select id from data_stores order by created_at limit 1);

DELETE FROM data_stores WHERE id != 'current';

COMMIT;
