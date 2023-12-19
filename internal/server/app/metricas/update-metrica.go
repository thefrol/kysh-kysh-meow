package metricas

import (
	"context"
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
)

func (mgr Manager) UpdateMetrica(ctx context.Context, m Metrica) (Metrica, error) {

	if m.ID == "" {
		return m, fmt.Errorf("%w: пустой айдишник", domain.ErrorBadID)
	}

	// получаем нужные данные из хранилища
	switch m.MType {
	case "counter":
		if m.Delta == nil {
			return m, fmt.Errorf("MetricaManager: %w: Delta==nil", domain.ErrorValidationError)
		}

		v, err := mgr.Registry.IncrementCounter(ctx, m.ID, *m.Delta)
		if err != nil {
			return m, fmt.Errorf("MetricaManager: %w", domain.ErrorMetricNotFound)
		}

		m.Delta = &v
	case "gauge":
		if m.Value == nil {
			return m, fmt.Errorf("MetricaManager: %w: Value==nil", domain.ErrorValidationError)
		}

		v, err := mgr.Registry.UpdateGauge(ctx, m.ID, *m.Value)
		if err != nil {
			return m, fmt.Errorf("MetricaManager: %w", domain.ErrorMetricNotFound)
		}

		m.Value = &v
	default:
		return m, fmt.Errorf("MetricaManager: %w %v ", domain.ErrorUnknownMetric, m.MType)
	}

	// возвращаем
	return m, nil
}
