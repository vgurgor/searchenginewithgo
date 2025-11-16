package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"search_engine/internal/domain/entities"
	"search_engine/internal/domain/repositories"
)

type contentMetricsRepository struct {
	pool *pgxpool.Pool
}

func NewContentMetricsRepository(pool *pgxpool.Pool) repositories.ContentMetricsRepository {
	return &contentMetricsRepository{pool: pool}
}

func (r *contentMetricsRepository) Create(ctx context.Context, m *entities.ContentMetrics) error {
	const q = `
		INSERT INTO content_metrics(content_id, views, likes, reading_time, reactions, final_score, recalculated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, q,
		m.ContentID, m.Views, m.Likes, m.ReadingTime, m.Reactions, m.FinalScore, m.RecalculatedAt,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *contentMetricsRepository) UpdateByContentID(ctx context.Context, contentID int64, m *entities.ContentMetrics) error {
	const q = `
		UPDATE content_metrics
		SET views=$1, likes=$2, reading_time=$3, reactions=$4, final_score=$5, recalculated_at=$6, updated_at=NOW()
		WHERE content_id=$7
		RETURNING id, updated_at
	`
	return r.pool.QueryRow(ctx, q,
		m.Views, m.Likes, m.ReadingTime, m.Reactions, m.FinalScore, m.RecalculatedAt, contentID,
	).Scan(&m.ID, &m.UpdatedAt)
}

func (r *contentMetricsRepository) GetByContentID(ctx context.Context, contentID int64) (*entities.ContentMetrics, error) {
	const q = `
		SELECT id, content_id, views, likes, reading_time, reactions, final_score, recalculated_at, created_at, updated_at
		FROM content_metrics WHERE content_id=$1
	`
	var m entities.ContentMetrics
	if err := r.pool.QueryRow(ctx, q, contentID).Scan(
		&m.ID, &m.ContentID, &m.Views, &m.Likes, &m.ReadingTime, &m.Reactions, &m.FinalScore, &m.RecalculatedAt, &m.CreatedAt, &m.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *contentMetricsRepository) BulkUpsert(ctx context.Context, metrics []entities.ContentMetrics) error {
	if len(metrics) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for i := range metrics {
		m := metrics[i]
		batch.Queue(`
			INSERT INTO content_metrics(content_id, views, likes, reading_time, reactions, final_score, recalculated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			ON CONFLICT (content_id) DO UPDATE SET
				views=EXCLUDED.views,
				likes=EXCLUDED.likes,
				reading_time=EXCLUDED.reading_time,
				reactions=EXCLUDED.reactions,
				final_score=EXCLUDED.final_score,
				recalculated_at=EXCLUDED.recalculated_at,
				updated_at=NOW()
		`, m.ContentID, m.Views, m.Likes, m.ReadingTime, m.Reactions, m.FinalScore, m.RecalculatedAt)
	}
	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	for range metrics {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}


