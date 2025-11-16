package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"search_engine/internal/domain/entities"
	"search_engine/internal/domain/repositories"
)

type contentRepository struct {
	pool *pgxpool.Pool
}

func NewContentRepository(pool *pgxpool.Pool) repositories.ContentRepository {
	return &contentRepository{pool: pool}
}

func (r *contentRepository) Create(ctx context.Context, c *entities.Content) error {
	const q = `
		INSERT INTO contents(
			provider_id, provider_content_id, title, content_type, description, url, thumbnail_url, published_at
		) VALUES($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, q,
		c.ProviderID, c.ProviderContentID, c.Title, c.ContentType, c.Description, c.URL, c.ThumbnailURL, c.PublishedAt,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *contentRepository) GetByID(ctx context.Context, id int64) (*entities.Content, error) {
	const q = `
		SELECT id, provider_id, provider_content_id, title, content_type, description, url, thumbnail_url, published_at, created_at, updated_at
		FROM contents WHERE id=$1
	`
	var c entities.Content
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&c.ID, &c.ProviderID, &c.ProviderContentID, &c.Title, &c.ContentType, &c.Description, &c.URL, &c.ThumbnailURL, &c.PublishedAt, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *contentRepository) GetAll(ctx context.Context, filters repositories.ContentFilters, pagination repositories.Pagination, sort repositories.SortBy) ([]entities.Content, int64, error) {
	base := `
		FROM contents c
		LEFT JOIN content_metrics cm ON cm.content_id = c.id
		WHERE 1=1
	`
	var args []any
	var where []string
	arg := 1
	if filters.ContentType != nil {
		where = append(where, fmt.Sprintf("c.content_type = $%d", arg))
		args = append(args, *filters.ContentType)
		arg++
	}
	if filters.ProviderID != nil {
		where = append(where, fmt.Sprintf("c.provider_id = $%d", arg))
		args = append(args, *filters.ProviderID)
		arg++
	}
	if len(where) > 0 {
		base += " AND " + strings.Join(where, " AND ")
	}
	countSQL := "SELECT COUNT(*) " + base
	var total int64
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	order := "ORDER BY c.published_at DESC"
	if sort == repositories.SortByPopularity {
		order = "ORDER BY cm.final_score DESC NULLS LAST"
	}

	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.PageSize <= 0 || pagination.PageSize > 100 {
		pagination.PageSize = 20
	}
	offset := (pagination.Page - 1) * pagination.PageSize

	selectSQL := `
		SELECT c.id, c.provider_id, c.provider_content_id, c.title, c.content_type, c.description, c.url, c.thumbnail_url, c.published_at, c.created_at, c.updated_at
	` + base + `
	` + order + `
	LIMIT $%d OFFSET $%d
	`
	selectSQL = fmt.Sprintf(selectSQL, arg, arg+1)
	args = append(args, pagination.PageSize, offset)

	rows, err := r.pool.Query(ctx, selectSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []entities.Content
	for rows.Next() {
		var c entities.Content
		if err := rows.Scan(
			&c.ID, &c.ProviderID, &c.ProviderContentID, &c.Title, &c.ContentType, &c.Description, &c.URL, &c.ThumbnailURL, &c.PublishedAt, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, c)
	}
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}
	return items, total, nil
}

func (r *contentRepository) Update(ctx context.Context, c *entities.Content) error {
	const q = `
		UPDATE contents
		SET provider_id=$1, provider_content_id=$2, title=$3, content_type=$4, description=$5, url=$6, thumbnail_url=$7, published_at=$8, updated_at=NOW()
		WHERE id=$9
		RETURNING updated_at
	`
	return r.pool.QueryRow(ctx, q,
		c.ProviderID, c.ProviderContentID, c.Title, c.ContentType, c.Description, c.URL, c.ThumbnailURL, c.PublishedAt, c.ID,
	).Scan(&c.UpdatedAt)
}

func (r *contentRepository) BulkInsert(ctx context.Context, contents []entities.Content) error {
	if len(contents) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for i := range contents {
		c := contents[i]
		batch.Queue(`
			INSERT INTO contents(provider_id, provider_content_id, title, content_type, description, url, thumbnail_url, published_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
			ON CONFLICT (provider_id, provider_content_id) DO NOTHING
		`, c.ProviderID, c.ProviderContentID, c.Title, c.ContentType, c.Description, c.URL, c.ThumbnailURL, c.PublishedAt)
	}
	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	for range contents {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}

func (r *contentRepository) SearchByKeyword(ctx context.Context, keyword string, filters repositories.ContentFilters, pagination repositories.Pagination, sort repositories.SortBy) ([]entities.Content, int64, error) {
	base := `
		FROM contents c
		LEFT JOIN content_metrics cm ON cm.content_id = c.id
		WHERE (c.title ILIKE $1 OR coalesce(c.description,'') ILIKE $1)
	`
	k := "%" + keyword + "%"
	args := []any{k}
	arg := 2
	if filters.ContentType != nil {
		base += fmt.Sprintf(" AND c.content_type = $%d", arg)
		args = append(args, *filters.ContentType)
		arg++
	}
	if filters.ProviderID != nil {
		base += fmt.Sprintf(" AND c.provider_id = $%d", arg)
		args = append(args, *filters.ProviderID)
		arg++
	}

	var total int64
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) "+base, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	order := "ORDER BY c.published_at DESC"
	if sort == repositories.SortByPopularity {
		order = "ORDER BY cm.final_score DESC NULLS LAST"
	}
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.PageSize <= 0 || pagination.PageSize > 100 {
		pagination.PageSize = 20
	}
	offset := (pagination.Page - 1) * pagination.PageSize

	sql := fmt.Sprintf(`
		SELECT c.id, c.provider_id, c.provider_content_id, c.title, c.content_type, c.description, c.url, c.thumbnail_url, c.published_at, c.created_at, c.updated_at
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, base, order, arg, arg+1)
	args = append(args, pagination.PageSize, offset)

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []entities.Content
	for rows.Next() {
		var c entities.Content
		if err := rows.Scan(
			&c.ID, &c.ProviderID, &c.ProviderContentID, &c.Title, &c.ContentType, &c.Description, &c.URL, &c.ThumbnailURL, &c.PublishedAt, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, c)
	}
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}
	return items, total, nil
}


