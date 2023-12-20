package metricas

import (
	"context"
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
)

func (mgr Manager) UpdateMetrica(ctx context.Context, m Metrica) (Metrica, error) {

	if m.ID == "" {
		return m, fmt.Errorf("%w: пустой айдишник", app.ErrorBadID)
	}

	// получаем нужные данные из хранилища
	switch m.MType {
	case "counter":
		if m.Delta == nil {
			return m, fmt.Errorf("MetricaManager: %w: Delta==nil", app.ErrorValidationError)
		}

		v, err := mgr.Registry.IncrementCounter(ctx, m.ID, *m.Delta)
		if err != nil {
			return m, fmt.Errorf("MetricaManager: %w", app.ErrorMetricNotFound)
		}

		m.Delta = &v
	case "gauge":
		if m.Value == nil {
			return m, fmt.Errorf("MetricaManager: %w: Value==nil", app.ErrorValidationError)
		}

		v, err := mgr.Registry.UpdateGauge(ctx, m.ID, *m.Value)
		if err != nil {
			return m, fmt.Errorf("MetricaManager: %w", app.ErrorMetricNotFound)
		}

		m.Value = &v
	default:
		return m, fmt.Errorf("MetricaManager: %w %v ", app.ErrorUnknownMetric, m.MType)
	}

	// возвращаем
	return m, nil
}
