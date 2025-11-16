package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"search_engine/internal/domain/entities"
	"search_engine/internal/domain/repositories"
	"strconv"
)

type syncHistoryRepository struct {
	pool *pgxpool.Pool
}

func NewSyncHistoryRepository(pool *pgxpool.Pool) repositories.SyncHistoryRepository {
	return &syncHistoryRepository{pool: pool}
}

func (r *syncHistoryRepository) Create(ctx context.Context, h *entities.SyncHistory) error {
	const q = `
		INSERT INTO sync_history(
			provider_id, sync_status, total_fetched, new_contents, updated_contents, skipped_contents, failed_contents, error_message, started_at, completed_at, duration_ms
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id
	`
	return r.pool.QueryRow(ctx, q,
		h.ProviderID, h.SyncStatus, h.TotalFetched, h.NewContents, h.UpdatedContents, h.SkippedContents, h.FailedContents, h.ErrorMessage, h.StartedAt, h.CompletedAt, h.DurationMs,
	).Scan(&h.ID)
}

func (r *syncHistoryRepository) Update(ctx context.Context, h *entities.SyncHistory) error {
	const q = `
		UPDATE sync_history
		SET sync_status=$1,
		    total_fetched=$2,
		    new_contents=$3,
		    updated_contents=$4,
		    skipped_contents=$5,
		    failed_contents=$6,
		    error_message=$7,
		    completed_at=$8,
		    duration_ms=$9
		WHERE id=$10
	`
	_, err := r.pool.Exec(ctx, q,
		h.SyncStatus,
		h.TotalFetched,
		h.NewContents,
		h.UpdatedContents,
		h.SkippedContents,
		h.FailedContents,
		h.ErrorMessage,
		h.CompletedAt,
		h.DurationMs,
		h.ID,
	)
	return err
}

func (r *syncHistoryRepository) GetByProviderID(ctx context.Context, providerID string, limit int) ([]entities.SyncHistory, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, provider_id, sync_status, total_fetched, new_contents, updated_contents, skipped_contents, failed_contents, error_message, started_at, completed_at, duration_ms
		FROM sync_history WHERE provider_id=$1 ORDER BY started_at DESC LIMIT $2
	`, providerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entities.SyncHistory
	for rows.Next() {
		var h entities.SyncHistory
		if err := rows.Scan(&h.ID, &h.ProviderID, &h.SyncStatus, &h.TotalFetched, &h.NewContents, &h.UpdatedContents, &h.SkippedContents, &h.FailedContents, &h.ErrorMessage, &h.StartedAt, &h.CompletedAt, &h.DurationMs); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, rows.Err()
}

func (r *syncHistoryRepository) GetLastSync(ctx context.Context, providerID string) (*entities.SyncHistory, error) {
	var h entities.SyncHistory
	err := r.pool.QueryRow(ctx, `
		SELECT id, provider_id, sync_status, total_fetched, new_contents, updated_contents, skipped_contents, failed_contents, error_message, started_at, completed_at, duration_ms
		FROM sync_history WHERE provider_id=$1 ORDER BY started_at DESC LIMIT 1
	`, providerID).Scan(&h.ID, &h.ProviderID, &h.SyncStatus, &h.TotalFetched, &h.NewContents, &h.UpdatedContents, &h.SkippedContents, &h.FailedContents, &h.ErrorMessage, &h.StartedAt, &h.CompletedAt, &h.DurationMs)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *syncHistoryRepository) GetAll(ctx context.Context, limit int) ([]entities.SyncHistory, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, provider_id, sync_status, total_fetched, new_contents, updated_contents, skipped_contents, failed_contents, error_message, started_at, completed_at, duration_ms
		FROM sync_history ORDER BY started_at DESC LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entities.SyncHistory
	for rows.Next() {
		var h entities.SyncHistory
		if err := rows.Scan(&h.ID, &h.ProviderID, &h.SyncStatus, &h.TotalFetched, &h.NewContents, &h.UpdatedContents, &h.SkippedContents, &h.FailedContents, &h.ErrorMessage, &h.StartedAt, &h.CompletedAt, &h.DurationMs); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, rows.Err()
}

func (r *syncHistoryRepository) List(ctx context.Context, providerID *string, status *entities.SyncStatus, limit, offset int) ([]entities.SyncHistory, error) {
	q := `
		SELECT id, provider_id, sync_status, total_fetched, new_contents, updated_contents, skipped_contents, failed_contents, error_message, started_at, completed_at, duration_ms
		FROM sync_history WHERE 1=1
	`
	args := []any{}
	arg := 1
	if providerID != nil && *providerID != "" {
		q += ` AND provider_id = $` + strconv.Itoa(arg)
		args = append(args, *providerID)
		arg++
	}
	if status != nil && *status != "" {
		q += ` AND sync_status = $` + strconv.Itoa(arg)
		args = append(args, *status)
		arg++
	}
	q += ` ORDER BY started_at DESC LIMIT $` + strconv.Itoa(arg) + ` OFFSET $` + strconv.Itoa(arg+1)
	args = append(args, limit, offset)
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entities.SyncHistory
	for rows.Next() {
		var h entities.SyncHistory
		if err := rows.Scan(&h.ID, &h.ProviderID, &h.SyncStatus, &h.TotalFetched, &h.NewContents, &h.UpdatedContents, &h.SkippedContents, &h.FailedContents, &h.ErrorMessage, &h.StartedAt, &h.CompletedAt, &h.DurationMs); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, rows.Err()
}

func (r *syncHistoryRepository) Count(ctx context.Context, providerID *string, status *entities.SyncStatus) (int64, error) {
	q := `SELECT COUNT(*) FROM sync_history WHERE 1=1`
	args := []any{}
	arg := 1
	if providerID != nil && *providerID != "" {
		q += ` AND provider_id = $` + strconv.Itoa(arg)
		args = append(args, *providerID)
		arg++
	}
	if status != nil && *status != "" {
		q += ` AND sync_status = $` + strconv.Itoa(arg)
		args = append(args, *status)
		arg++
	}
	var total int64
	if err := r.pool.QueryRow(ctx, q, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}


