package providers

import (
	"fmt"
	domainp "search_engine/internal/domain/providers"
)

type ProviderFactory struct {
	registry map[string]domainp.IContentProvider
}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{registry: make(map[string]domainp.IContentProvider)}
}

func (f *ProviderFactory) RegisterProvider(p domainp.IContentProvider) {
	f.registry[p.GetProviderID()] = p
}

func (f *ProviderFactory) GetAllProviders() []domainp.IContentProvider {
	list := make([]domainp.IContentProvider, 0, len(f.registry))
	for _, p := range f.registry {
		list = append(list, p)
	}
	return list
}

func (f *ProviderFactory) GetProviderByID(id string) (domainp.IContentProvider, error) {
	p, ok := f.registry[id]
	if !ok {
		return nil, fmt.Errorf("provider not found: %s", id)
	}
	return p, nil
}
