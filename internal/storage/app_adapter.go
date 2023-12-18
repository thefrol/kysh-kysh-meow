package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/manager"
	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
)

var (
	_ manager.CounterRepository = (*CounterAdapter)(nil)
	_ manager.GaugeRepository   = (*GaugeAdapter)(nil)
)

type CounterAdapter struct {
	Op api.Operator
}

// Get implements manager.CounterRepository.
func (adapter *CounterAdapter) Get(ctx context.Context, id string) (int64, error) {
	// это запрос по сути
	d := metrica.Metrica{
		MType: "counter",
		ID:    id,
	}

	// получаем метрику из оператора
	v, err := adapter.Op.Get(ctx, d)
	if err != nil {
		if errors.Is(err, api.ErrorNotFoundMetric) {
			return 0, fmt.Errorf("in CounterAdapter: %w: %v", domain.ErrorMetricNotFound, err)
		}
		return 0, err
	}

	// проверяем, что мы что-то нормальное получили)
	if len(v) != 1 {
		return 0, fmt.Errorf("in COunterAdapter: получено невалидное значение из оператора, больше одного значения или ноль")
	}

	if v[0].Delta == nil {
		return 0, fmt.Errorf("in COunterAdapter: получено невалидное значение из оператора, ссылка на metrica.Delta==nil")
	}

	return *v[0].Delta, nil
}

// Increment implements manager.CounterRepository.
func (adapter *CounterAdapter) Increment(ctx context.Context, id string, delta int64) (int64, error) {
	// это запрос по сути
	d := metrica.Metrica{
		MType: "counter",
		ID:    id,
		Delta: &delta,
	}

	// получаем метрику из оператора
	v, err := adapter.Op.Update(ctx, d)
	if err != nil {
		if errors.Is(err, api.ErrorNotFoundMetric) {
			return 0, fmt.Errorf("in CounterAdapter: %w: %v", domain.ErrorMetricNotFound, err)
		}
		return 0, err
	}

	// проверяем, что мы что-то нормальное получили)
	if len(v) != 1 {
		return 0, fmt.Errorf("in COunterAdapter: получено невалидное значение из оператора, больше одного значения или ноль")
	}

	if v[0].Delta == nil {
		return 0, fmt.Errorf("in COunterAdapter: получено невалидное значение из оператора, ссылка на metrica.Delta==nil")
	}

	return *v[0].Delta, nil
}

type GaugeAdapter struct {
	Op api.Operator
}

// Get implements manager.GaugeRepository.
func (adapter *GaugeAdapter) Get(ctx context.Context, id string) (float64, error) {
	// это запрос по сути
	d := metrica.Metrica{
		MType: "gauge",
		ID:    id,
	}

	// получаем метрику из оператора
	v, err := adapter.Op.Get(ctx, d)
	if err != nil {
		if errors.Is(err, api.ErrorNotFoundMetric) {
			return 0, fmt.Errorf("in GaugeAdapter: %w: %v", domain.ErrorMetricNotFound, err)
		}
		return 0, err
	}

	// проверяем, что мы что-то нормальное получили)
	if len(v) != 1 {
		return 0, fmt.Errorf("in GaugeAdapter: получено невалидное значение из оператора, больше одного значения или ноль")
	}

	if v[0].Value == nil {
		return 0, fmt.Errorf("in GaugeAdapter: получено невалидное значение из оператора, ссылка на metrica.Delta==nil")
	}

	return *v[0].Value, nil
}

// Increment implements manager.GaugeRepository.
func (adapter *GaugeAdapter) Update(ctx context.Context, id string, value float64) (float64, error) {
	// это запрос по сути
	d := metrica.Metrica{
		MType: "gauge",
		ID:    id,
		Value: &value,
	}

	// получаем метрику из оператора
	v, err := adapter.Op.Update(ctx, d)
	if err != nil {
		if errors.Is(err, api.ErrorNotFoundMetric) {
			return 0, fmt.Errorf("in GaugeAdapter: %w: %v", domain.ErrorMetricNotFound, err)
		}
		return 0, err
	}

	// проверяем, что мы что-то нормальное получили)
	if len(v) != 1 {
		return 0, fmt.Errorf("in GaugeAdapter: получено невалидное значение из оператора, больше одного значения или ноль")
	}

	if v[0].Value == nil {
		return 0, fmt.Errorf("in GaugeAdapter: получено невалидное значение из оператора, ссылка на metrica.Delta==nil")
	}

	return *v[0].Value, nil
}