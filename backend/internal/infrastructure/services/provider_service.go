package services

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	domainp "search_engine/internal/domain/providers"
	"search_engine/internal/infrastructure/ratelimiter"
)

type ProviderService struct {
	Factory interface {
		GetAllProviders() []domainp.IContentProvider
		GetProviderByID(id string) (domainp.IContentProvider, error)
	}
	Limiter *ratelimiter.RedisLimiter
	Logger  *zap.Logger
	Timeout time.Duration
}

func (s *ProviderService) FetchFromAllProviders(ctx context.Context) ([]domainp.ProviderContent, error) {
	providers := s.Factory.GetAllProviders()
	var wg sync.WaitGroup
	type result struct {
		items []domainp.ProviderContent
		err   error
	}
	resCh := make(chan result, len(providers))
	for _, p := range providers {
		p := p
		wg.Add(1)
		go func() {
			defer wg.Done()
			providerID := p.GetProviderID()
			ok, err := s.Limiter.CheckLimit(ctx, providerID, p.GetRateLimit().RequestsPerMinute)
			if err != nil {
				s.Logger.Warn("rate limit check error", zap.String("provider", providerID), zap.Error(err))
			}
			if !ok {
				resCh <- result{nil, nil}
				return
			}
			_ = s.Limiter.RecordRequest(ctx, providerID)
			cctx, cancel := context.WithTimeout(ctx, s.Timeout)
			defer cancel()
			done := make(chan result, 1)
			go func() {
				items, e := p.FetchContents()
				done <- result{items, e}
			}()
			select {
			case r := <-done:
				if r.err != nil {
					s.Logger.Warn("provider fetch failed", zap.String("provider", providerID), zap.Error(r.err))
				}
				resCh <- r
			case <-cctx.Done():
				s.Logger.Warn("provider fetch timeout", zap.String("provider", providerID))
				resCh <- result{nil, cctx.Err()}
			}
		}()
	}
	wg.Wait()
	close(resCh)
	var all []domainp.ProviderContent
	for r := range resCh {
		if len(r.items) > 0 {
			all = append(all, r.items...)
		}
	}
	return all, nil
}

func (s *ProviderService) FetchFromProvider(ctx context.Context, providerID string) ([]domainp.ProviderContent, error) {
	p, err := s.Factory.GetProviderByID(providerID)
	if err != nil {
		return nil, err
	}
	ok, err := s.Limiter.CheckLimit(ctx, providerID, p.GetRateLimit().RequestsPerMinute)
	if err != nil {
		s.Logger.Warn("rate limit check error", zap.String("provider", providerID), zap.Error(err))
	}
	if !ok {
		return []domainp.ProviderContent{}, nil
	}
	_ = s.Limiter.RecordRequest(ctx, providerID)
	cctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()
	done := make(chan struct {
		items []domainp.ProviderContent
		err   error
	}, 1)
	go func() {
		items, e := p.FetchContents()
		done <- struct {
			items []domainp.ProviderContent
			err   error
		}{items, e}
	}()
	select {
	case r := <-done:
		return r.items, r.err
	case <-cctx.Done():
		return nil, cctx.Err()
	}
}
