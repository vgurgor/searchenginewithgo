package providers

import "testing"

func TestProviderFactory_RegisterAndGet(t *testing.T) {
	f := NewProviderFactory()
	p := NewJSONProvider("http://example", 0)
	f.RegisterProvider(p)
	got, err := f.GetProviderByID("provider1")
	if err != nil || got == nil {
		t.Fatalf("expected provider1, got err=%v", err)
	}
	all := f.GetAllProviders()
	if len(all) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(all))
	}
}
