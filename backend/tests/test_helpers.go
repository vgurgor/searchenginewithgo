package tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	rediscontainer "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

// TestEnvironment manages test database and redis connections
type TestEnvironment struct {
	DB      *pgxpool.Pool
	Redis   *redis.Client
	Logger  *zap.Logger
	Cleanup func()
}

// SetupTestEnvironment creates a test environment with PostgreSQL and Redis
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	t.Helper()

	ctx := context.Background()
	var cleanupFuncs []func()

	// Start PostgreSQL container
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}
	cleanupFuncs = append(cleanupFuncs, func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate PostgreSQL container: %v", err)
		}
	})

	// Get PostgreSQL connection string
	pgConnStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get PostgreSQL connection string: %v", err)
	}

	// Start Redis container
	redisContainer, err := rediscontainer.RunContainer(ctx,
		testcontainers.WithImage("redis:7-alpine"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithOccurrence(1).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("Failed to start Redis container: %v", err)
	}
	cleanupFuncs = append(cleanupFuncs, func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate Redis container: %v", err)
		}
	})

	// Get Redis connection string
	redisConnStr, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get Redis connection string: %v", err)
	}

	// Remove redis:// prefix if present
	redisAddr := strings.TrimPrefix(redisConnStr, "redis://")

	// Connect to PostgreSQL
	dbPool, err := pgxpool.New(ctx, pgConnStr)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Test database connection
	if err := dbPool.Ping(ctx); err != nil {
		t.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Test Redis connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to ping Redis: %v", err)
	}

	// Run migrations
	if err := runMigrations(ctx, pgConnStr); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create test logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	cleanup := func() {
		for i := len(cleanupFuncs) - 1; i >= 0; i-- {
			cleanupFuncs[i]()
		}
		dbPool.Close()
		_ = rdb.Close()
		_ = logger.Sync()
	}

	return &TestEnvironment{
		DB:      dbPool,
		Redis:   rdb,
		Logger:  logger,
		Cleanup: cleanup,
	}
}

// runMigrations runs database migrations for testing
func runMigrations(ctx context.Context, connStr string) error {
	// For now, we'll use a simple approach to run migrations
	// In a real scenario, you might want to use golang-migrate or similar
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func() { _ = db.Close() }()

	// Run basic migrations for testing
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS pg_trgm;`,
		`DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'content_type_enum') THEN
				CREATE TYPE content_type_enum AS ENUM ('video', 'text');
			END IF;
		END $$;`,
		`CREATE TABLE IF NOT EXISTS contents (
			id BIGSERIAL PRIMARY KEY,
			provider_id VARCHAR(50) NOT NULL,
			provider_content_id VARCHAR(100) NOT NULL,
			title VARCHAR(255) NOT NULL,
			content_type content_type_enum NOT NULL,
			description TEXT NULL,
			url VARCHAR(500) NULL,
			thumbnail_url VARCHAR(500) NULL,
			published_at TIMESTAMPTZ NULL,
			deleted_at TIMESTAMPTZ NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(provider_id, provider_content_id)
		);`,
		`CREATE TABLE IF NOT EXISTS content_metrics (
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
		);`,
		`CREATE INDEX IF NOT EXISTS idx_contents_lower_title ON contents((LOWER(title)));
		CREATE INDEX IF NOT EXISTS idx_contents_lower_description ON contents((LOWER(description)));
		CREATE INDEX IF NOT EXISTS idx_contents_published_at ON contents(published_at DESC);
		CREATE INDEX IF NOT EXISTS idx_content_metrics_final_score_desc ON content_metrics(final_score DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_contents_fts ON contents USING GIN (to_tsvector('english', title || ' ' || COALESCE(description, '')));`,
		`CREATE OR REPLACE FUNCTION content_search_relevance(
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
			title_score := ts_rank(to_tsvector('english', title), plainto_tsquery('english', search_query)) * title_weight;
			desc_score := ts_rank(to_tsvector('english', COALESCE(description, '')), plainto_tsquery('english', search_query)) * desc_weight;
			RETURN title_score + desc_score;
		END;
		$$ LANGUAGE plpgsql IMMUTABLE;`,
		`CREATE OR REPLACE FUNCTION fuzzy_search_relevance(
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
			title_similarity := similarity(search_term, title) * title_weight;
			desc_similarity := similarity(search_term, COALESCE(description, '')) * desc_weight;
			RETURN title_similarity + desc_similarity;
		END;
		$$ LANGUAGE plpgsql IMMUTABLE;`,
	}

	for _, migration := range migrations {
		if _, err := db.ExecContext(ctx, migration); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	return nil
}

// SkipIfShort skips the test if running in short mode
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}
}

// SkipIfNoDocker skips the test if Docker is not available
func SkipIfNoDocker(t *testing.T) {
	if os.Getenv("NO_DOCKER") == "true" {
		t.Skip("Skipping test requiring Docker")
	}
}
