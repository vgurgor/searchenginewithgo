-- Full-text search indexes and functions
-- Create GIN index for full-text search
CREATE INDEX IF NOT EXISTS idx_contents_fts ON contents USING GIN (to_tsvector('english', title || ' ' || COALESCE(description, '')));

-- Create function for relevance ranking
CREATE OR REPLACE FUNCTION content_search_relevance(
    search_query TEXT,
    title TEXT,
    description TEXT
) RETURNS REAL AS $$
DECLARE
    title_weight REAL := 1.0;
    desc_weight REAL := 0.4;
    title_score REAL;
    desc_score REAL;
BEGIN
    -- Calculate title relevance
    title_score := ts_rank(to_tsvector('english', title), plainto_tsquery('english', search_query)) * title_weight;

    -- Calculate description relevance
    desc_score := ts_rank(to_tsvector('english', COALESCE(description, '')), plainto_tsquery('english', search_query)) * desc_weight;

    -- Return combined score
    RETURN title_score + desc_score;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Create function for fuzzy search with trigrams
CREATE OR REPLACE FUNCTION fuzzy_search_relevance(
    search_term TEXT,
    title TEXT,
    description TEXT
) RETURNS REAL AS $$
DECLARE
    title_similarity REAL;
    desc_similarity REAL;
    title_weight REAL := 0.7;
    desc_weight REAL := 0.3;
BEGIN
    -- Calculate trigram similarity for title
    title_similarity := similarity(search_term, title) * title_weight;

    -- Calculate trigram similarity for description
    desc_similarity := similarity(search_term, COALESCE(description, '')) * desc_weight;

    -- Return combined similarity score
    RETURN title_similarity + desc_similarity;
END;
$$ LANGUAGE plpgsql IMMUTABLE;
