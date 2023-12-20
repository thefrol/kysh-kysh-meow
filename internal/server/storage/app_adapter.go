package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/manager"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/scan"
	"github.com/thefrol/kysh-kysh-meow/internal/server/router/httpio"
)

var (
	_ manager.CounterRepository = (*CounterAdapter)(nil)
	_ scan.CounterLister        = (*CounterAdapter)(nil)
	_ manager.GaugeRepository   = (*GaugeAdapter)(nil)
	_ scan.GaugeLister          = (*GaugeAdapter)(nil)
)

type CounterAdapter struct {
	Op httpio.Operator
}

// All implements scan.CounterLister.
func (adapter *CounterAdapter) All(ctx context.Context) (map[string]int64, error) {
	res := make(map[string]int64)

	counters, _, err := adapter.Op.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, l := range counters {
		v, err := adapter.Counter(ctx, l)
		if err != nil {
			return nil, err
		}
		res[l] = v
	}

	return res, nil
}

// Counter implements manager.CounterRepository.
func (adapter *CounterAdapter) Counter(ctx context.Context, id string) (int64, error) {
	// это запрос по сути
	d := metrica.Metrica{
		MType: "counter",
		ID:    id,
	}

	// получаем метрику из оператора
	v, err := adapter.Op.Get(ctx, d)
	if err != nil {
		if errors.Is(err, httpio.ErrorNotFoundMetric) {
			return 0, fmt.Errorf("in CounterAdapter: %w: %v", app.ErrorMetricNotFound, err)
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

// CounterIncrement implements manager.CounterRepository.
func (adapter *CounterAdapter) CounterIncrement(ctx context.Context, id string, delta int64) (int64, error) {
	// это запрос по сути
	d := metrica.Metrica{
		MType: "counter",
		ID:    id,
		Delta: &delta,
	}

	// получаем метрику из оператора
	v, err := adapter.Op.Update(ctx, d)
	if err != nil {
		if errors.Is(err, httpio.ErrorNotFoundMetric) {
			return 0, fmt.Errorf("in CounterAdapter: %w: %v", app.ErrorMetricNotFound, err)
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
	Op httpio.Operator
}

// All implements scan.GaugeLister.
func (adapter *GaugeAdapter) All(ctx context.Context) (map[string]float64, error) {
	res := make(map[string]float64)

	_, gauges, err := adapter.Op.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, l := range gauges {
		v, err := adapter.Gauge(ctx, l)
		if err != nil {
			return nil, err
		}
		res[l] = v
	}

	return res, nil
}

// Gauge implements manager.GaugeRepository.
func (adapter *GaugeAdapter) Gauge(ctx context.Context, id string) (float64, error) {
	// это запрос по сути
	d := metrica.Metrica{
		MType: "gauge",
		ID:    id,
	}

	// получаем метрику из оператора
	v, err := adapter.Op.Get(ctx, d)
	if err != nil {
		if errors.Is(err, httpio.ErrorNotFoundMetric) {
			return 0, fmt.Errorf("in GaugeAdapter: %w: %v", app.ErrorMetricNotFound, err)
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
func (adapter *GaugeAdapter) GaugeUpdate(ctx context.Context, id string, value float64) (float64, error) {
	// это запрос по сути
	d := metrica.Metrica{
		MType: "gauge",
		ID:    id,
		Value: &value,
	}

	// получаем метрику из оператора
	v, err := adapter.Op.Update(ctx, d)
	if err != nil {
		if errors.Is(err, httpio.ErrorNotFoundMetric) {
			return 0, fmt.Errorf("in GaugeAdapter: %w: %v", app.ErrorMetricNotFound, err)
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
