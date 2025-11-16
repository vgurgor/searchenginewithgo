package services

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"

	"search_engine/internal/domain/entities"
	"search_engine/internal/domain/providers"
	"search_engine/internal/domain/repositories"
)

type fakeFactory struct{ items []providers.ProviderContent }

func (f *fakeFactory) GetAllProviders() []providers.IContentProvider {
	return []providers.IContentProvider{&fakeProvider{items: f.items}}
}
func (f *fakeFactory) GetProviderByID(id string) (providers.IContentProvider, error) {
	return &fakeProvider{items: f.items}, nil
}

type fakeProvider struct{ items []providers.ProviderContent }

func (p *fakeProvider) FetchContents() ([]providers.ProviderContent, error) { return p.items, nil }
func (p *fakeProvider) GetProviderID() string                               { return "provider1" }
func (p *fakeProvider) GetRateLimit() providers.RateLimit {
	return providers.RateLimit{RequestsPerMinute: 100}
}

type fakeProviderClient struct{ items []providers.ProviderContent }

func (f *fakeProviderClient) FetchFromProvider(ctx context.Context, providerID string) ([]providers.ProviderContent, error) {
	return f.items, nil
}

type memContentRepo struct {
	byKey map[string]*entities.Content
	all   []*entities.Content
}

func (m *memContentRepo) key(pid, cid string) string { return pid + "|" + cid }
func (m *memContentRepo) Create(ctx context.Context, c *entities.Content) error {
	if m.byKey == nil {
		m.byKey = map[string]*entities.Content{}
	}
	c.ID = int64(len(m.all) + 1)
	now := time.Now().UTC()
	c.CreatedAt, c.UpdatedAt = now, now
	m.byKey[m.key(c.ProviderID, c.ProviderContentID)] = c
	m.all = append(m.all, c)
	return nil
}
func (m *memContentRepo) GetByProviderKey(ctx context.Context, providerID, providerContentID string) (*entities.Content, error) {
	if x, ok := m.byKey[m.key(providerID, providerContentID)]; ok {
		return x, nil
	}
	return nil, ErrNotFound
}
func (m *memContentRepo) GetByID(ctx context.Context, id int64) (*entities.Content, error) {
	for _, c := range m.all {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, ErrNotFound
}
func (m *memContentRepo) Update(ctx context.Context, c *entities.Content) error {
	c.UpdatedAt = time.Now().UTC()
	return nil
}
func (m *memContentRepo) GetAll(ctx context.Context, _ repositories.ContentFilters, _ repositories.Pagination, _ repositories.SortBy) ([]entities.Content, int64, error) {
	return nil, 0, nil
}
func (m *memContentRepo) BulkInsert(ctx context.Context, contents []entities.Content) error {
	return nil
}
func (m *memContentRepo) SearchByKeyword(ctx context.Context, keyword string, filters repositories.ContentFilters, pagination repositories.Pagination, sort repositories.SortBy) ([]entities.Content, int64, error) {
	return nil, 0, nil
}
func (m *memContentRepo) ListIDs(ctx context.Context, offset, limit int) ([]int64, error) {
	return nil, nil
}
func (m *memContentRepo) CountAll(ctx context.Context) (int64, error) { return int64(len(m.all)), nil }
func (m *memContentRepo) SearchWithFilters(ctx context.Context, keyword string, contentType *entities.ContentType, pagination repositories.Pagination, sort repositories.SearchSort) ([]repositories.ContentWithMetrics, int64, error) {
	return nil, 0, nil
}
func (m *memContentRepo) GetDetailByID(ctx context.Context, id int64) (*repositories.ContentWithMetrics, error) {
	return nil, ErrNotFound
}
func (m *memContentRepo) CountByType(ctx context.Context) (map[entities.ContentType]int64, error) {
	return map[entities.ContentType]int64{}, nil
}
func (m *memContentRepo) GetAverageScore(ctx context.Context) (float64, error) { return 0, nil }
func (m *memContentRepo) CountByProvider(ctx context.Context) (map[string]int64, error) {
	return map[string]int64{}, nil
}
func (m *memContentRepo) SoftDelete(ctx context.Context, id int64) error { return nil }
func (m *memContentRepo) ListIDsByType(ctx context.Context, t entities.ContentType, offset, limit int) ([]int64, error) {
	return nil, nil
}
func (m *memContentRepo) GetAverageScoreByProvider(ctx context.Context, providerID string) (float64, error) {
	return 0, nil
}

type Err struct{}

var ErrNotFound = &Err{}

func (*Err) Error() string { return "not found" }

type memMetricsRepo struct {
	byID map[int64]*entities.ContentMetrics
}

func (r *memMetricsRepo) Create(ctx context.Context, m *entities.ContentMetrics) error {
	if r.byID == nil {
		r.byID = map[int64]*entities.ContentMetrics{}
	}
	m.ID = int64(len(r.byID) + 1)
	now := time.Now().UTC()
	m.CreatedAt, m.UpdatedAt = now, now
	r.byID[m.ContentID] = m
	return nil
}
func (r *memMetricsRepo) UpdateByContentID(ctx context.Context, contentID int64, m *entities.ContentMetrics) error {
	if r.byID == nil {
		r.byID = map[int64]*entities.ContentMetrics{}
	}
	old := r.byID[contentID]
	if old == nil {
		r.byID[contentID] = m
		return nil
	}
	*old = *m
	return nil
}
func (r *memMetricsRepo) GetByContentID(ctx context.Context, contentID int64) (*entities.ContentMetrics, error) {
	if x, ok := r.byID[contentID]; ok {
		return x, nil
	}
	return nil, ErrNotFound
}
func (r *memMetricsRepo) BulkUpsert(ctx context.Context, metrics []entities.ContentMetrics) error {
	return nil
}

type noopHistoryRepo struct{}

func (n *noopHistoryRepo) Create(ctx context.Context, h *entities.SyncHistory) error {
	h.ID++
	return nil
}
func (n *noopHistoryRepo) Update(ctx context.Context, h *entities.SyncHistory) error { return nil }
func (n *noopHistoryRepo) GetByProviderID(ctx context.Context, providerID string, limit int) ([]entities.SyncHistory, error) {
	return nil, nil
}
func (n *noopHistoryRepo) GetLastSync(ctx context.Context, providerID string) (*entities.SyncHistory, error) {
	return nil, nil
}
func (n *noopHistoryRepo) GetAll(ctx context.Context, limit int) ([]entities.SyncHistory, error) {
	return nil, nil
}
func (n *noopHistoryRepo) List(ctx context.Context, providerID *string, status *entities.SyncStatus, limit, offset int) ([]entities.SyncHistory, error) {
	return nil, nil
}
func (n *noopHistoryRepo) Count(ctx context.Context, providerID *string, status *entities.SyncStatus) (int64, error) {
	return 0, nil
}

func TestContentSyncService_NewAndUpdate(t *testing.T) {
	logger := zap.NewNop() // Use no-op logger for tests
	now := time.Now().UTC()
	items := []providers.ProviderContent{
		{ProviderID: "provider1", ProviderContentID: "a1", Title: "T1", ContentType: "text", PublishedAt: now},
		{ProviderID: "provider1", ProviderContentID: "v1", Title: "V1", ContentType: "video", PublishedAt: now},
	}
	factory := &fakeFactory{items: items}
	crepo := &memContentRepo{}
	mrepo := &memMetricsRepo{}
	engine := &mockEngine{}
	scoreCalc := &ScoreCalculatorService{Contents: crepo, Metrics: mrepo, Engine: engine, Logger: logger}
	providerClient := &fakeProviderClient{items: items}
	svc := &ContentSyncService{
		Logger:         logger,
		Factory:        factory,
		ProviderClient: providerClient,
		Contents:       crepo,
		Metrics:        mrepo,
		ScoreCalc:      scoreCalc,
		HistoryRepo:    &noopHistoryRepo{},
		Thresholds:     MetricsThresholds{Percent: 5, AbsViews: 100, AbsLikes: 10, AbsReactions: 5},
	}
	ctx := context.Background()
	res, err := svc.SyncProvider(ctx, "provider1")
	if err != nil {
		t.Fatalf("sync error: %v", err)
	}
	if res.NewContents != 2 {
		t.Fatalf("expected 2 new contents, got %d", res.NewContents)
	}
	// update metrics path
	items2 := []providers.ProviderContent{
		{ProviderID: "provider1", ProviderContentID: "a1", Title: "T1", ContentType: "text", PublishedAt: now, Reactions: intPtr(100)},
	}
	svc.Factory = &fakeFactory{items: items2}
	svc.ProviderClient = &fakeProviderClient{items: items2}
	res2, _ := svc.SyncProvider(ctx, "provider1")
	if res2.UpdatedContents == 0 {
		t.Fatalf("expected updated contents > 0")
	}
}

type mockEngine struct{}

func (m *mockEngine) CalculateScore(content entities.Content, metrics entities.ContentMetrics) (float64, error) {
	return 42.0, nil
}

func intPtr(v int) *int { return &v }
