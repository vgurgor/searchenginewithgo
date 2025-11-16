-- Enable pg_trgm extension for fast ILIKE searches
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- GIN trigram indexes to accelerate ILIKE searches on title and description
CREATE INDEX IF NOT EXISTS idx_contents_title_trgm ON contents USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_contents_description_trgm ON contents USING gin (description gin_trgm_ops);



