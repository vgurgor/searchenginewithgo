-- Enum for content type
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'content_type_enum') THEN
        CREATE TYPE content_type_enum AS ENUM ('video', 'text');
    END IF;
END $$;

-- contents table
CREATE TABLE IF NOT EXISTS contents (
    id BIGSERIAL PRIMARY KEY,
    provider_id VARCHAR(50) NOT NULL,
    provider_content_id VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content_type content_type_enum NOT NULL,
    description TEXT NULL,
    url VARCHAR(500) NULL,
    thumbnail_url VARCHAR(500) NULL,
    published_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(provider_id, provider_content_id)
);

-- content_metrics table
CREATE TABLE IF NOT EXISTS content_metrics (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGINT NOT NULL REFERENCES contents(id) ON DELETE CASCADE ON UPDATE CASCADE,
    views BIGINT NOT NULL DEFAULT 0,
    likes BIGINT NOT NULL DEFAULT 0,
    reading_time INT NOT NULL DEFAULT 0,
    reactions INT NOT NULL DEFAULT 0,
    final_score NUMERIC(10,2) NOT NULL DEFAULT 0,
    recalculated_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(content_id)
);

-- indexes for contents
CREATE INDEX IF NOT EXISTS idx_contents_provider_id ON contents(provider_id);
CREATE INDEX IF NOT EXISTS idx_contents_content_type ON contents(content_type);
CREATE INDEX IF NOT EXISTS idx_contents_published_at ON contents(published_at DESC);

-- indexes for content_metrics
CREATE INDEX IF NOT EXISTS idx_content_metrics_content_id ON content_metrics(content_id);
CREATE INDEX IF NOT EXISTS idx_content_metrics_final_score_desc ON content_metrics(final_score DESC);

-- NOTE: Cross-table composite indexes are not supported in PostgreSQL.
-- The intended composite index (content_type + final_score DESC) is represented via
-- separate indexes on contents(content_type) and content_metrics(final_score DESC).


