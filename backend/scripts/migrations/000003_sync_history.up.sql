DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'sync_status_enum') THEN
        CREATE TYPE sync_status_enum AS ENUM ('success','partial','failed','in_progress','skipped');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS sync_history (
    id BIGSERIAL PRIMARY KEY,
    provider_id VARCHAR(50) NOT NULL,
    sync_status sync_status_enum NOT NULL,
    total_fetched INT NOT NULL DEFAULT 0,
    new_contents INT NOT NULL DEFAULT 0,
    updated_contents INT NOT NULL DEFAULT 0,
    skipped_contents INT NOT NULL DEFAULT 0,
    failed_contents INT NOT NULL DEFAULT 0,
    error_message TEXT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ NULL,
    duration_ms INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_sync_history_provider_started_at ON sync_history(provider_id, started_at DESC);


