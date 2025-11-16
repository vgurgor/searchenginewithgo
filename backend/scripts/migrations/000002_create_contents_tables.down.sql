DROP TABLE IF EXISTS content_metrics;
DROP TABLE IF EXISTS contents;
DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'content_type_enum') THEN
        DROP TYPE content_type_enum;
    END IF;
END $$;


