-- Case-insensitive LIKE search support
CREATE INDEX IF NOT EXISTS idx_contents_lower_title ON contents((LOWER(title)));
CREATE INDEX IF NOT EXISTS idx_contents_lower_description ON contents((LOWER(description)));
-- Published date already has index, ensure exists
CREATE INDEX IF NOT EXISTS idx_contents_published_at2 ON contents(published_at DESC);


