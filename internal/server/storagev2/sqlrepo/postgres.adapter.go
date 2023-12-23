package sqlrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/manager"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/scan"
)

var (
	_ manager.CounterRepository = (*Repository)(nil)
	_ manager.GaugeRepository   = (*Repository)(nil)
	_ scan.Labler               = (*Repository)(nil)
)

var (
	ErrorNilQueries = fmt.Errorf("адаптер sqlc==nil: %w", app.ErrorNilReference)
)

// Repository это класс-адаптер между
// Queries(запросами в бд) и юз-кейсами
// прикладного уровня
//
// Из запросов в БД генерируется класс queries,
// но там своеобразные сигнатуры и методы,
// и возвращаемые ошибки, мы хотим привести это к некоему стандарту
type Repository struct {
	Q *Queries

	Log zerolog.Logger
}

// Labels implements scan.Labler.
func (repo *Repository) Labels(ctx context.Context) (map[string][]string, error) {
	if repo == nil {
		return nil, fmt.Errorf("postgres.adapter: %w", app.ErrorNilReference)
	}

	if repo.Q == nil {
		return nil, fmt.Errorf("postgres.adapter: %w", ErrorNilQueries)
	}

	l, err := repo.Q.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgres.adapter: %w", err)
	}

	// Теперь конвертируем в наш формат вывода

	// массив имен счетчиков
	var cs []string
	var gs []string

	for _, v := range l {
		switch v.Column1 {
		case "counter":
			cs = append(cs, v.ID)
		case "gauge":
			cs = append(gs, v.ID)
		default:
			repo.Log.Error().
				Str("operation", "List").
				Str("id", v.ID).
				Str("type", v.Column1).
				Msg("неизвестный тип метрики")
		}
	}

	m := map[string][]string{
		"counters": cs,
		"gauges":   gs,
	}

	return m, nil
}

// Gauge implements manager.GaugeRepository.
func (repo *Repository) Gauge(ctx context.Context, id string) (float64, error) {
	if repo == nil {
		return 0, fmt.Errorf("postgres.adapter: %w", app.ErrorNilReference)
	}

	if repo.Q == nil {
		return 0, fmt.Errorf("postgres.adapter: %w", ErrorNilQueries)
	}

	c, err := repo.Q.Gauge(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("postgres.adapter: %w", app.ErrorMetricNotFound)
		}
		return 0, fmt.Errorf("postgres.adapter: %w", err)
	}

	return c.Value, nil
}

// GaugeUpdate implements manager.GaugeRepository.
func (repo *Repository) GaugeUpdate(ctx context.Context, id string, value float64) (float64, error) {
	if repo == nil {
		return 0, fmt.Errorf("postgres.adapter: %w", app.ErrorNilReference)
	}

	if repo.Q == nil {
		return 0, fmt.Errorf("postgres.adapter: %w", ErrorNilQueries)
	}

	v := UpdateGaugeParams{
		ID:    id,
		Value: value,
	}

	c, err := repo.Q.UpdateGauge(ctx, v)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("postgres.adapter: %w", app.ErrorMetricNotFound)
		}
		return 0, fmt.Errorf("postgres.adapter: %w", err)
	}

	return c.Value, nil
}

// Counter implements manager.CounterRepository.
func (repo *Repository) Counter(ctx context.Context, id string) (int64, error) {
	if repo == nil {
		return 0, fmt.Errorf("postgres.adapter: %w", app.ErrorNilReference)
	}

	if repo.Q == nil {
		return 0, fmt.Errorf("postgres.adapter: %w", ErrorNilQueries)
	}

	c, err := repo.Q.Counter(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("postgres.adapter: %w", app.ErrorMetricNotFound)
		}
		return 0, fmt.Errorf("postgres.adapter: %w", err)
	}

	return c.Value, nil
}

// CounterIncrement implements manager.CounterRepository.
func (repo *Repository) CounterIncrement(ctx context.Context, id string, delta int64) (int64, error) {
	if repo == nil {
		return 0, fmt.Errorf("postgres.adapter: %w", app.ErrorNilReference)
	}

	if repo.Q == nil {
		return 0, fmt.Errorf("postgres.adapter: %w", ErrorNilQueries)
	}

	inc := IncrementCounterParams{
		ID:    id,
		Value: delta,
	}

	c, err := repo.Q.IncrementCounter(ctx, inc)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("postgres.adapter: %w", app.ErrorMetricNotFound)
		}
		return 0, fmt.Errorf("postgres.adapter: %w", err)
	}

	return c.Value, nil
}
