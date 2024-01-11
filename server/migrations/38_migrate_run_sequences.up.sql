DO $$
DECLARE
    temprow record;
BEGIN
    FOR temprow IN
        SELECT max(id)+1 AS next_value, test_id, tenant_id
        FROM test_runs
        GROUP BY test_id, tenant_id
    LOOP
        EXECUTE format('CREATE SEQUENCE IF NOT EXISTS runs_test_%s_seq START WITH %s',
                       MD5(FORMAT('%s%s', temprow.test_id, temprow.tenant_id)), temprow.next_value);
    END LOOP;

    FOR temprow IN
        SELECT max(id::int)+1 AS next_value, test_suite_id, tenant_id
        FROM test_suite_runs
        GROUP BY test_suite_id, tenant_id
    LOOP
        EXECUTE format('CREATE SEQUENCE IF NOT EXISTS runs_test_suite_%s_seq START WITH %s',
                       MD5(FORMAT('%s%s', temprow.test_suite_id, temprow.tenant_id)), temprow.next_value);
    END LOOP;
END $$;
