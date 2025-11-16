DROP TABLE IF EXISTS sync_history;
DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'sync_status_enum') THEN
        DROP TYPE sync_status_enum;
    END IF;
END $$;


