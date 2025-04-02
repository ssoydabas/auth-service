#!/bin/bash

# Create test database
PGPASSWORD=postgres psql -h localhost -U postgres -c "DROP DATABASE IF EXISTS auth_service_test;"
PGPASSWORD=postgres psql -h localhost -U postgres -c "CREATE DATABASE auth_service_test;"

# Run migrations on test database
PGPASSWORD=postgres psql -h localhost -U postgres -d auth_service_test -c "
DO $$
DECLARE
    table_name text;
BEGIN
    FOR table_name IN
        SELECT tablename
        FROM pg_tables
        WHERE schemaname = 'public'
    LOOP
        EXECUTE 'DROP TABLE IF EXISTS ' || table_name || ' CASCADE';
    END LOOP;
END $$;
"

# Run the tests
go test -v -count=1 ./tests/integration/... 