-- Drop full-text search indexes and functions
DROP INDEX IF EXISTS idx_contents_fts;
DROP FUNCTION IF EXISTS content_search_relevance(TEXT, TEXT, TEXT);
DROP FUNCTION IF EXISTS fuzzy_search_relevance(TEXT, TEXT, TEXT);
