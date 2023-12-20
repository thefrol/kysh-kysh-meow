package metricas

import (
	"context"
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
)

func (mgr Manager) GetMetrica(ctx context.Context, m Metrica) (Metrica, error) {

	// сначала валидируем входные данные
	err := ValidateRequest(m)
	if err != nil {
		return m, fmt.Errorf("MetricaManager: Metrica не прошла валидацию %w", err)
	}

	// получаем нужные данные из хранилища
	switch m.MType {
	case "counter":
		v, err := mgr.Registry.Counter(ctx, m.ID)
		if err != nil {
			return m, fmt.Errorf("MetricaManager: %w", app.ErrorMetricNotFound)
		}

		m.Delta = &v
	case "gauge":
		v, err := mgr.Registry.Gauge(ctx, m.ID)
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
