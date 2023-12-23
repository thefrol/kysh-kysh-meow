package manager

import "context"

type CounterRepository interface {
	Counter(ctx context.Context, id string) (int64, error)
	CounterIncrement(ctx context.Context, id string, delta int64) (int64, error)
}

type GaugeRepository interface {
	Gauge(ctx context.Context, id string) (float64, error)
	GaugeUpdate(ctx context.Context, id string, value float64) (float64, error)
}

// Registry это реестр метрик,
// он позволяет сохранять и удалять
// метрики. Это некая абстракция над хранилищем.
type Registry struct {
	Counters CounterRepository
	Gauges   GaugeRepository
}
