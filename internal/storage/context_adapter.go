package storage

import (
	"context"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
)

// ContextAdapter оборачивает классы на старом апи хранилища oldAPI,
// под новое апи api.Storager, делая так, контекст, конечно используется тупо вхолостую
type ContextAdapter struct {
	legacyStore Storager
}

// NewAdapter оборачивает хранилища старого интерфейса, позволяя их подключать к хендлерам, работающим на новом интерфейсе api.Storager
func NewAdapter(s Storager) *ContextAdapter {
	return &ContextAdapter{legacyStore: s}
}

// Counter implements api.Storager.
func (a ContextAdapter) Counter(ctx context.Context, name string) (value int64, err error) {
	c, found := a.legacyStore.Counter(name)
	if !found {
		return 0, api.ErrorNotFoundMetric
	}
	return int64(c), nil
}

// Gauge implements api.Storager.
func (a ContextAdapter) Gauge(ctx context.Context, name string) (value float64, err error) {
	g, found := a.legacyStore.Gauge(name)
	if !found {
		return 0, api.ErrorNotFoundMetric
	}
	return float64(g), nil
}

// IncrementCounter implements api.Storager.
func (a *ContextAdapter) IncrementCounter(ctx context.Context, name string, delta int64) (value int64, err error) {
	was, _ := a.legacyStore.Counter(name)
	newVal := was + metrica.Counter(delta)
	a.legacyStore.SetCounter(name, newVal)
	return int64(newVal), nil
}

// UpdateGauge implements api.Storager.
func (a *ContextAdapter) UpdateGauge(ctx context.Context, name string, v float64) (value float64, err error) {
	a.legacyStore.SetGauge(name, metrica.Gauge(v))
	return v, nil
}

// List implements api.Storager.
func (a *ContextAdapter) List(ctx context.Context) (counterNames []string, gaugeNames []string, err error) {
	return a.legacyStore.ListCounters(), a.legacyStore.ListGauges(), nil
}

var _ api.Storager = (*ContextAdapter)(nil)
